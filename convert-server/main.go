package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"
)

// LambdaJSON is request parameter
type LambdaJSON struct {
	Language string `json:"language"`
	Version  string `json:"version"`
	Program  string `json:"program"`
}

var reg, _ = regexp.Compile("language:\\s?(.*)\nversion:\\s?(.*)\n```\n?([\\s\\S]*)\n?```")

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		errorLog("invalid method")
		w.WriteHeader(405)
		fmt.Fprint(w, "invalid method")
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	log("request: \n" + string(body))

	m := reg.FindSubmatch(body)

	if len(m) != 4 {
		errorLog("invalid request parameter")
		w.WriteHeader(400)
		fmt.Fprint(w, "invalid request parameter")
		return
	}

	json, e := json.Marshal(LambdaJSON{string(m[1]), string(m[2]), string(m[3])})

	if e != nil {
		errorLog("invalid request parameter")
		w.WriteHeader(400)
		fmt.Fprint(w, "invalid request parameter")
		return
	}

	req, err := http.NewRequest("POST", "http://localhost:8081", bytes.NewBuffer(json))
	if err != nil {
		errorLog("invalid request parameter")
		w.WriteHeader(400)
		fmt.Fprint(w, "invalid request parameter")
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	resbody, _ := ioutil.ReadAll(resp.Body)

	if err != nil {
		errorLog("lambda server error: " + string(resbody))
		w.WriteHeader(resp.StatusCode)
		fmt.Fprint(w, "lambda server error: "+string(resbody))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	log("response: " + string(resbody))
	fmt.Fprint(w, string(resbody))
}

func log(msg string) {
	fmt.Printf("[%s] %s\n", time.Now(), msg)
}

func errorLog(msg string) {
	fmt.Fprintf(os.Stderr, "[%s] %s\n", time.Now(), msg)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
