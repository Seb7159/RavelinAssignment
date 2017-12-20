package main

import (
	"encoding/json"
	"fmt" 
	"log"
	"net/http"
)

type helloWorldRequest struct { 
	Hello string 
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) { 
	decoder := json.NewDecoder(r.Body) 

	var req helloWorldRequest 
	var err = decoder.Decode(&req); 
	
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to decode JSON request"))
		//return
	}

	defer r.Body.Close() 
	log.Printf("Request received %+v", req)

	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/", helloWorldHandler)

	fmt.Println("Server now running on localhost:8080")
	fmt.Println(`Try: curl -X POST -d "{\"hello\": \"that\"}" http://localhost:8080`)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
