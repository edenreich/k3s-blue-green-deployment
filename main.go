
package main

import (
	"log"
	"net/http"
	"fmt"
	"os"
)

func main()  {
	fmt.Println("Started server on 0.0.0.0:80")
	calls := 0
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request)  {
		calls++

		hostname, _ := os.Hostname();
		
		fmt.Printf("Recieved request nr. %d\n", calls)
		w.Write([]byte(fmt.Sprintf("hello from %s on version 1", hostname)))
	})

	log.Fatal(http.ListenAndServe("0.0.0.0:80", nil))
}