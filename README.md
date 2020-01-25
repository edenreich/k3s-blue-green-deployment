# Blue-Green-Deployment

Short example of how to apply blue green deployment to k3s raspberry pi cluster.

By blue green deployment is normally meant we are going to double our resources (so for example if we have 4 pods, we are going to have 8 pods in total during this process). 2 services will be available with two different versions of our application. Our loadbalancer will direct traffic to the service that exposes the deployment of application version 1 (see `blue.yaml`), while application version 1 is running, we will deploy internally another service that points to application version 2(see `green.yaml`) and eventually we'll configure our loadbalancer to move the traffic to our green deployment and cleanup the blue one. 

## Requirements

This example assume you have 4-5 Raspberry pis. 1 as master and the rest are workers.
Installed version: `v1.17.0+k3s.1`

Private docker registry is needed to host the docker images.

When I run `kubectl get nodes`, I get the following output:

```sh
k4s-master   Ready    master   4d22h   v1.17.0+k3s.1
k4s-node1    Ready    worker   4d21h   v1.17.0+k3s.1
k4s-node2    Ready    worker   4d21h   v1.17.0+k3s.1
k4s-node3    Ready    worker   4d21h   v1.17.0+k3s.1
k4s-node4    Ready    worker   4d21h   v1.17.0+k3s.1
```

You should've already setup your own cluster, there are plenty of tutorials explaining how to do this.

## Walkthrough

First we need to create a secret out of our docker configuration file for out private docker registry so the cluster can pull and run docker our images, run:

```sh
docker login

kubectl create secret generic regcred \
    --from-file=.dockerconfigjson=$HOME/.docker/config.json \
    --type=kubernetes.io/dockerconfigjson
```

On a separate terminal session lets run an infinity curl with a sleep of 1 sec that will hit our loadbalancer, so we can simulate a customer constantly hiting our website:

```sh 
while true; do curl http://test.website.com; echo -e; sleep 1; done
```

`test.website.com` is a local domain name that resolves to the ip address of my loadbalancer, to figure out what's the ip address of your loadbalancer, run `kubectl describe ingress`, on loadbalancer.yaml I've also setup that the hostname of the master node resolves to the ip address of the loadbalancer, so you could simply configure the hostname of your master node in `loadbalancer.yaml` and use that hostname for the curl command, for simplicity purposes this loadbalancer directing any traffic regardless of what hostname is requested (*) to the given service name on a given port.

Next we need to deploy our loadbalancer (k3s comes with build-in `traefik` controller, let's use that):

```sh
kubectl apply -f loadbalancer.yaml
```

Let build application version 1:

1. Build version 1 hello-world app: `docker build -t [account]/test:1 .` (where account is the name of your account, you might need to login first - just run `docker login`).
2. Push it to our registry: `docker push`.
3. Deploy version 1: `kubectl apply -f blue.yaml`.

Now we want to upgrade to application version 2 without our customers to experince down-times, so we need to prepare our next version before we tell our loadbalancer to direct all the traffic to our new version.

first, let's modify the output of our application to `version 2` so we can notice the changes from our curl session.

1. Build version 2 hello-world app: `docker build -t [account]/test:2 .`.
2. Push it to our registry: `docker push`.
3. Deploy version 2: `kubectl apply -f green.yaml`.

So we deployed another service that points to 4 replicas of our application version 2. We can verifiy it by running: `kubectl get pods -w`. You should see that instead of 4 replicas we have now 8 pods. 4 pods under each service version.

Even thought we deployed our green service that points to deployment version 2, we still didn't tell our loadbalancer to point the traffic to our new deployment, so we need to tell it to switch the traffic to `web-app-2` in service name section.
Once done modifing let's apply the changes: `kubectl apply -f loadbalancer.yaml`.
You should now see that our curl client session is being served from our new version.
After we are 100% sure this version is good enough, we need to delete the previous deployment to save some resources, run: `kubectl delete -f blue.yaml`.

