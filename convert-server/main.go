package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
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
		w.WriteHeader(405)
		fmt.Fprint(w, "invalid method")
		return
	}

	body, _ := ioutil.ReadAll(r.Body)

	m := reg.FindSubmatch(body)

	if len(m) != 4 {
		w.WriteHeader(400)
		fmt.Fprint(w, "invalid request parameter")
		return
	}

	json, e := json.Marshal(LambdaJSON{string(m[1]), string(m[2]), string(m[3])})

	if e != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "invalid request parameter")
		return
	}

	req, err := http.NewRequest("POST", "http://localhost:8081", bytes.NewBuffer(json))
	if err != nil {
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
		w.WriteHeader(resp.StatusCode)
		fmt.Fprint(w, "lambda server error: "+string(resbody))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(resbody))
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
