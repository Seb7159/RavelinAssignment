package main

import (
	"encoding/json"
	"fmt" 
	"log"
	"net/http"
	"strings"
)


var i	 int  = 0; 
var data map[string]Data; 


type Data struct {
	WebsiteUrl         string
	SessionId          string
	ResizeFrom         Dimension
	ResizeTo           Dimension
	CopyAndPaste       map[string]bool 
	FormCompletionTime int 
}

type Dimension struct {
	Width  string
	Height string
} 


func findId(sessId string) string {
	for index, element := range data {
		if strings.Compare(sessId, element.SessionId) == 0 {
			return index 
		}
	}
	return "error"
}


func dataHandler(w http.ResponseWriter, r *http.Request) { 
	decoder := json.NewDecoder(r.Body) 

	var dt  Data 
	var err = decoder.Decode(&dt); 
	
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to decode JSON request"))
		return
	}
	defer r.Body.Close() 

	if dt.ResizeTo.Width == ""	{
		log.Println("COPY-PASTE detected from session ID: " + dt.SessionId)
		log.Println(dt.CopyAndPaste)
		}  else {
		log.Println("The data was COMPLETED by the following ID: " + dt.SessionId)
		log.Println("The whole struct is: ")  
		log.Println(dt) 
		} 

	data[dt.SessionId] = dt; 

	w.WriteHeader(http.StatusOK) 
}



func main() { 
	data = make(map[string]Data); 
	data["initialise"] = Data{}; 
	
	// Handle the main page 
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:]) 
	}) 

	// Handle JSON requests 
	http.HandleFunc("/data", dataHandler) 

	// Show message when server is up 
	fmt.Println("Server now running on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
