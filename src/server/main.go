package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Profile struct {
	Name    string
	Hobbies []string
}

type Response struct {
	Url         string
	Paramaters  string
	RequestType string
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello world</h1>")
}

func about(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>About us</h1>")
}

func customAPI(w http.ResponseWriter, r *http.Request) {

	var responseData map[string]interface{}
	switch r.Method {
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		requestType := r.FormValue("request_type")
		paramters, _ := json.Marshal(r.FormValue("parameters"))
		url := r.FormValue("url")
		switch requestType {
		case "get":
			resp, err := http.Get(url)
			if err != nil {
				// handle error
			}
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			json.Unmarshal(body, &responseData)
		case "post":
			fmt.Println("making post request at " + url + " with parameters " + r.FormValue("parameters"))
			resp, err := http.Post(url, "application/json", bytes.NewBuffer(paramters))

			if err != nil {
				fmt.Printf("ERror occured: ")
				fmt.Println(err)
			}
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			println(body)
			json.Unmarshal(body, &responseData)
		}
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}

	js, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/about", about)
	http.HandleFunc("/custom_api", customAPI)
	fmt.Println("Server starting...")
	http.ListenAndServe(":3000", nil)
}
