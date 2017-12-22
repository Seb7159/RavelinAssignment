package main

import (
	"encoding/json"
	"fmt" 
	"log"
	"net/http" 
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


func dataHandler(w http.ResponseWriter, r *http.Request) { 
	// Decode the JSON 
	decoder := json.NewDecoder(r.Body) 

	var dt  Data 
	var err = decoder.Decode(&dt); 
	
	// Handle errors 
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to decode JSON request"))
		return
	}
	defer r.Body.Close() 

	// Print data in the console 
	fmt.Println("");
	if dt.ResizeTo.Width == ""	{
		log.Println("COPY-PASTE detected from session ID: " + dt.SessionId)
		log.Println(dt.CopyAndPaste)
		}  else {
		log.Println("The data was COMPLETED by the following ID: " + dt.SessionId)
		log.Println("The whole struct is: ")  
		log.Println(dt) 
		} 

	// Put JSON in the 'data' map 
	data[dt.SessionId] = dt; 

	// Send status OK in HTTP 
	w.WriteHeader(http.StatusOK) 
}



func main() { 
	// Initialise map of 'data' 
	data = make(map[string]Data); 
	
	// Handle the main page 
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:]) 
	}) 

	// Handle JSON requests 
	http.HandleFunc("/data", dataHandler) 

	// OPTIONAL: Display data map by accessing the link: /show
	http.HandleFunc("/show", func(w http.ResponseWriter, r *http.Request) { 
		fmt.Println("");
		fmt.Println("");
		fmt.Println(""); 

		fmt.Println("The whole 'data' map is: "); 
		for index := range data{ 
			fmt.Println(data[index]) 
		}

		fmt.Println("");
		fmt.Println("");
		fmt.Fprintf(w, "See the terminal! "); 
	}) 

	// Show message when server is up 
	fmt.Println("Server now running on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
