package main

import (
	"encoding/json"
	"fmt" 
	"log"
	"net/http" 
	"sync" 
)


// Declare the 'data' map variable globally 
var data sync.Map; 


// Temporary JSON request struct 
type tempJSONrequest struct {
	EventType		   string 			`json:"eventType"`
	FormId			   string 			`json:"formId"`
	WebsiteUrl         string 			`json:"websiteUrl"`
	SessionId          string			`json:"sessionId"`
	ResizeFrom         Dimension		
	ResizeTo           Dimension
	CopyAndPaste       map[string]bool 	
	Pasted			   bool 			`json:"pasted"`
	FormCompletionTime int 				`json:"time"` 
}


// Default given structs TODO: add name json to each of those 
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
	var dt  tempJSONrequest  
	var err = decoder.Decode(&dt) 
	
	// Handle errors 
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to decode JSON request"))
		return
	}
	defer r.Body.Close() 

	// Make a temporary Data variable 
	var tempData Data 
	// If the key was initialised before 
	if result, ok := data.Load(dt.SessionId); ok {
	    tempData.WebsiteUrl 			= result.(Data).WebsiteUrl
	    tempData.SessionId 				= result.(Data).SessionId
	    tempData.ResizeFrom			 	= result.(Data).ResizeFrom
	    tempData.ResizeTo  				= result.(Data).ResizeTo
	    tempData.CopyAndPaste  			= result.(Data).CopyAndPaste
	    tempData.FormCompletionTime  	= result.(Data).FormCompletionTime 
	// Else if it was not 
	} else { tempData = Data{} } 
	
	// Print data in the console 
	log.Println(); 
	// In case there was copy-paste detected 
	if dt.EventType == "copyAndPaste"	{
		log.Println("COPY-PASTE detected from session ID: " + dt.SessionId) 

		// Assign values 
		tempData.WebsiteUrl    = dt.WebsiteUrl 
		tempData.SessionId	   = dt.SessionId

		// Check if map for copyAndPaste was initialised before 
		if len(tempData.CopyAndPaste) == 0 {
					tempData.CopyAndPaste    = make(map[string]bool) 
		} 
		tempData.CopyAndPaste[dt.FormId]     = dt.Pasted 
		
		// Assign the temporary value to the map
		data.Store(dt.SessionId, tempData) 

		// Print the confirmation 
		log.Println("The following input was pasted: " + dt.FormId) 
		

	} else if dt.EventType == "resizeWindow" {
		// In case the window was resized 
		log.Println("RESIZE detected from the following ID: " + dt.SessionId) 

		// Assign values 
		tempData.WebsiteUrl 			   = dt.WebsiteUrl
		tempData.SessionId 				   = dt.SessionId 
		tempData.ResizeFrom.Width 		   = dt.ResizeFrom.Width
		tempData.ResizeFrom.Height 		   = dt.ResizeFrom.Height 
		tempData.ResizeTo.Width 		   = dt.ResizeTo.Width
		tempData.ResizeTo.Height		   = dt.ResizeTo.Height 


		// Assign the temporary value to the map
		data.Store(dt.SessionId, tempData) 

		// Print element struct 
		log.Println(tempData) 


	} else if dt.EventType == "timeTaken" { 
		// In case the submit button was pressed 
		log.Println("The struct was COMPLETED by the following ID: " + dt.SessionId) 

		// In case there were no copyAndPaste events before 
		if len(tempData.CopyAndPaste) == 0 { 
				// Initialise element if copyAndPaste or resize not happened 
				if tempData.ResizeFrom.Width == ""{
					tempData = Data{} 
				}
				tempData.CopyAndPaste    = make(map[string]bool) 
		} 
		
		// Assign values 
		tempData.WebsiteUrl 			   = dt.WebsiteUrl
		tempData.SessionId 				   = dt.SessionId 
		tempData.FormCompletionTime		   = dt.FormCompletionTime  
		
		// Check if the copyAndPaste map fields exist and initalise them 
		if ok := tempData.CopyAndPaste["inputEmail"]; !ok {
		    tempData.CopyAndPaste["inputEmail"] = false 
		}	
		if ok := tempData.CopyAndPaste["inputCardNumber"]; !ok {
		    tempData.CopyAndPaste["inputCardNumber"] = false 
		}	
		if ok := tempData.CopyAndPaste["inputCVV"]; !ok { 
		    tempData.CopyAndPaste["inputCVV"] = false 
		} 
		
		// Assign the temporary value to the map
		data.Store(dt.SessionId, tempData) 

		// Print element struct 
		log.Println(tempData)
		log.Println()  


	} else { // In case the eventType is not recognised 
		log.Println("ERROR! JSON request event type is not recognised! ")
		return; 
	}


	// Send status OK in HTTP 
	w.WriteHeader(http.StatusOK) 
}


// Print in console 'data' map method 
func showMap(w http.ResponseWriter, r *http.Request) {
		// Print map content 
		log.Println("The whole 'data' map is: ") 
		data.Range(func(key, value interface{}) bool {
			log.Println(value); 
			return true
		}) 
		log.Println("")
		log.Println("") 

		// Show message on the browser 
		fmt.Fprintf(w, "See the terminal! ") 
}


// Main method 
func main() { 
	// Handle index.html 
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:]) 
	}) 

	// Handle JSON requests 
	http.HandleFunc("/data", dataHandler) 

	// OPTIONAL: Display data map by accessing the link: /show
	http.HandleFunc("/show", showMap)  

	// Show message when server is up and run it on the 8080 port 
	log.Println("Server is now running on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
} 