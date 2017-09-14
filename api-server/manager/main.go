package manager

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"../docker-cli"
)

// In is request body
type In struct {
	Language string `json:"language"`
	Version  string `json:"version"`
	Program  string `json:"program"`
}

// Out is response body
type Out struct {
	Stdout string `json:"stdout"`
	Error  string `json:"error"`
}

// ContainerHandler execute job
// request body is ...
// POST
// {
//    "language": "Ruby",
//    "version": "2.3.0",
//    "code": "...",
// }
func ContainerHandler(w http.ResponseWriter, r *http.Request) {
	out := Out{}

	defer func() {
		outjson, _ := json.Marshal(out)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(outjson))
	}()

	if r.Method != "POST" {
		out.Error = "can use only post method"
		return
	}

	body, e := ioutil.ReadAll(r.Body)
	if e != nil {
		out.Error = e.Error()
		return
	}

	in := In{}
	e = json.Unmarshal(body, &in)
	if e != nil {
		out.Error = e.Error()
		return
	}

	// pipeline
	//
	//  build image
	//       |
	// run container(& remove)
	//       |
	// recept stdout
	image := dockerCli.Build(in.Language, in.Version, in.Program)
	stdout := dockerCli.Run(image)
	out.Stdout = <-stdout
}
