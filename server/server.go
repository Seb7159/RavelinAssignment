// Package name is set to 'main' 
package main

// Import all the libraries relevant to the server 
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
// Any JSON will be handled in this struct and the relevant data will be sent in the actual saved array of structs 
type tempJSONrequest struct {
	EventType		   string 			`json:"eventType"`
	FormId			   string 			`json:"formId"`
	WebsiteUrl         string 			`json:"websiteUrl"`
	SessionId          string			`json:"sessionId"`
	ResizeFrom         Dimension		
	ResizeTo           Dimension
	CopyAndPaste       sync.Map 	
	Pasted			   bool 			`json:"pasted"`
	FormCompletionTime int 				`json:"time"` 
}


// Default given data struct 
// Added sync map in order to make the server syncronized and thread-safe 
type Data struct {
	WebsiteUrl         string 
	SessionId          string
	ResizeFrom         Dimension
	ResizeTo           Dimension
	CopyAndPaste       sync.Map 
	FormCompletionTime int 
}

// Default struct for the dimension type 
type Dimension struct {
	Width  string
	Height string
} 


// Print completed struct method
// @param d   The data struct that has to be printed 
func printComplete(d Data){
		log.Println("Website URL: " + d.WebsiteUrl)
		log.Println("Session ID:  " + d.SessionId)

	if( d.ResizeFrom.Width != "" ){
		log.Println("ResizeFrom width:   " + d.ResizeFrom.Width)
		log.Println("ResizeFrom height:  " + d.ResizeFrom.Height)
		log.Println("ResizeTo width:     " + d.ResizeTo.Width)
		log.Println("ResizeTo height:    " + d.ResizeTo.Height)
	} else {
		log.Println("Resize event not detected") 
	}

		boolMail, ok := d.CopyAndPaste.Load("inputEmail")
		boolCardNo, _ := d.CopyAndPaste.Load("inputCardNumber")
		boolCVV, _    := d.CopyAndPaste.Load("inputCVV")

	// In case the struct was completed 
	if( ok == true ){
		log.Println("Copy and paste by form IDs: ")
		log.Println("        inputEmail:       ", boolMail ) 
		log.Println("        inputCardNumber:  ", boolCardNo )
		log.Println("        inputCVV:         ", boolCVV )
		log.Println("Form completion time:     ", d.FormCompletionTime) 

	// In case it is not complete 
	} else {
		log.Println("The struct was NOT completed. ") 
	}

		log.Println()  
		log.Println()  
}

// Data handler method 
// @param w    The response writer for the http
// @param r    The request from the http server
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

		// Store the data 
		tempData.CopyAndPaste.Store(dt.FormId, dt.Pasted)  
		
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
		log.Println("From (" + tempData.ResizeFrom.Width + ", " + tempData.ResizeFrom.Height +
					") to (" + tempData.ResizeTo.Width + ", " + tempData.ResizeTo.Height + ")") 


	} else if dt.EventType == "timeTaken" { 
		// In case the submit button was pressed 
		log.Println("The struct was COMPLETED by the following ID: " + dt.SessionId) 

		// Count how many fields were pasted 
		len := 0
		tempData.CopyAndPaste.Range(func(key, value interface{}) bool {
			len++ 
			return true 
		}) 

		// If not initialised, set to empty 
		if tempData.ResizeFrom.Width == "" && len == 0 {
			tempData = Data{} 
		} 
		
		// Assign values 
		tempData.WebsiteUrl 			   = dt.WebsiteUrl
		tempData.SessionId 				   = dt.SessionId 
		tempData.FormCompletionTime		   = dt.FormCompletionTime  
		
		if _, ok := tempData.CopyAndPaste.Load("inputEmail"); !ok {
		    tempData.CopyAndPaste.Store("inputEmail", false)  
		}	
		if _, ok := tempData.CopyAndPaste.Load("inputCardNumber"); !ok {
		    tempData.CopyAndPaste.Store("inputCardNumber", false)
		}	
		if _, ok := tempData.CopyAndPaste.Load("inputCVV"); !ok { 
		    tempData.CopyAndPaste.Store("inputCVV", false)  
		} 
		
		// Assign the temporary value to the map
		data.Store(dt.SessionId, tempData) 

		// Print element struct 
		printComplete(tempData) 


	} else { // In case the eventType is not recognised 
		log.Println("ERROR! JSON request event type is not recognised! ")
		return; 
	}


	// Send status OK in HTTP 
	w.WriteHeader(http.StatusOK) 
}


// Print in console 'data' map method 
// @param w    The response writer for the http
// @param r    The request from the http server 
func showMap(w http.ResponseWriter, r *http.Request) {
		// Count elements
		leng := 0 
		data.Range(func(key, value interface{}) bool {
			leng++ 
			return true
		}) 

		// Print map content depending whether there are elements or not 
		if leng != 0 { 
			log.Println()
			log.Println() 
			log.Println("The whole 'data' map is: ") 
			log.Println() 

			// For loop 
			data.Range(func(key, value interface{}) bool {
			printComplete(value.(Data)) 
			return true
			}) 

			// Display two empty lines for the design 
			log.Println("")
			log.Println("") 
		} else { 
			log.Println("The 'data' map is empty at the moment. ") 
		} 
		

		// Show message on the browser 
		fmt.Fprintf(w, "Check the terminal! ") 
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