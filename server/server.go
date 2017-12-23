package main

import (
	"encoding/json"
	"fmt" 
	"log"
	"net/http" 
)


// Declare the 'data' map variable globally 
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


// Data handler method 
func dataHandler(w http.ResponseWriter, r *http.Request) { 
	// Decode the JSON 
	decoder := json.NewDecoder(r.Body) 

	// Declare the temporary data and error variables 
	var dt  Data 
	var err = decoder.Decode(&dt) 
	
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
		log.Println("The struct was COMPLETED by the following ID: " + dt.SessionId)
		log.Println(dt) 
		} 

	// Put JSON in the 'data' map 
	data[dt.SessionId] = dt 

	// Send status OK in HTTP 
	w.WriteHeader(http.StatusOK) 
}


// Print in console 'data' map method 
func showMap(w http.ResponseWriter, r *http.Request) {
		// Print map content 
		fmt.Println("\n\n\nThe whole 'data' map is: ") 
		for index := range data{ 
			fmt.Println(data[index]) 
		} 
		fmt.Println("\n")  

		// Show message on the browser 
		fmt.Fprintf(w, "See the terminal! ") 
}


// Main method 
func main() { 
	// Initialise 'data' map 
	data = make(map[string]Data) 
	
	// Handle index.html 
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:]) 
	}) 

	// Handle JSON requests 
	http.HandleFunc("/data", dataHandler) 

	// OPTIONAL: Display data map by accessing the link: /show
	http.HandleFunc("/show", showMap)  

	// Show message when server is up and run it on the 8080 port 
	fmt.Println("Server now running on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
} 