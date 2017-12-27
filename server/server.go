package main

import (
	"encoding/json"
	"fmt" 
	"log"
	"net/http" 
)


// Declare the 'data' map variable globally 
var data map[string]*Data; 


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


	// Print data in the console 
	fmt.Println(); 
	// In case there was copy-paste detected 
	if dt.EventType == "copyAndPaste"	{
		log.Println("COPY-PASTE detected from session ID: " + dt.SessionId)
		
		// If the key was not initialised before 
		if _, ok := data[dt.SessionId]; !ok {
		    data[dt.SessionId] = &Data{} 
		} 

		// Assign values 
		data[dt.SessionId].WebsiteUrl 				   = dt.WebsiteUrl
		data[dt.SessionId].SessionId 				   = dt.SessionId

		// Check if map for copyAndPaste was initialised before 
		if len(data[dt.SessionId].CopyAndPaste) == 0 {
					data[dt.SessionId].CopyAndPaste    = make(map[string]bool) 
				} 
		data[dt.SessionId].CopyAndPaste[dt.FormId]     = dt.Pasted 
		
		// Print the confirmation 
		log.Println("The following input was pasted: " + dt.FormId) 
		

	} else if dt.EventType == "resizeWindow" {
		// In case the window was resized 
		log.Println("RESIZE detected from the following ID: " + dt.SessionId)

		// If the key was not initialised before 
		if _, ok := data[dt.SessionId]; !ok { 
		    data[dt.SessionId] = &Data{} 
		}

		// Assign values 
		data[dt.SessionId].ResizeFrom.Width 		   = dt.ResizeFrom.Width
		data[dt.SessionId].ResizeFrom.Height 		   = dt.ResizeFrom.Height 
		data[dt.SessionId].ResizeTo.Width 			   = dt.ResizeTo.Width
		data[dt.SessionId].ResizeTo.Height			   = dt.ResizeTo.Height 

		// Print element struct 
		log.Println(*data[dt.SessionId]) 


	} else if dt.EventType == "timeTaken" { 
		// In case the submit button was pressed 
		log.Println("The struct was COMPLETED by the following ID: " + dt.SessionId)
		
		// If the key was not initialised before 
		if _, ok := data[dt.SessionId]; !ok {
		    data[dt.SessionId] = &Data{} 
		}
		
		// In case there were no copyAndPaste events before 
		if len(data[dt.SessionId].CopyAndPaste) == 0 { 
				// Initialise element if copyAndPaste or resize not happened 
				if data[dt.SessionId].ResizeFrom.Width == ""{
					data[dt.SessionId] = &Data{} 
				}
				data[dt.SessionId].CopyAndPaste    = make(map[string]bool) 
		} 
		
		// Assign values 
		data[dt.SessionId].WebsiteUrl 				   = dt.WebsiteUrl
		data[dt.SessionId].SessionId 				   = dt.SessionId 
		data[dt.SessionId].FormCompletionTime		   = dt.FormCompletionTime  
		
		// Check if the copyAndPaste map fields exist and initalise them 
		if ok := data[dt.SessionId].CopyAndPaste["inputEmail"]; !ok {
		    data[dt.SessionId].CopyAndPaste["inputEmail"] = false 
		}	
		if ok := data[dt.SessionId].CopyAndPaste["inputCardNumber"]; !ok {
		    data[dt.SessionId].CopyAndPaste["inputCardNumber"] = false 
		}	
		if ok := data[dt.SessionId].CopyAndPaste["inputCVV"]; !ok { 
		    data[dt.SessionId].CopyAndPaste["inputCVV"] = false 
		} 
		
		// Print element struct 
		log.Println(*data[dt.SessionId]) 


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
	data = make(map[string]*Data) 
	
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