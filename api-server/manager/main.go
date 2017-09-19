package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

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
		log("response: " + string(outjson))
		fmt.Fprint(w, string(outjson))
	}()

	if r.Method != "POST" {
		w.WriteHeader(405)
		errorLog("invalid method")
		fmt.Fprint(w, "invalid method")
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	log("response: " + string(body))

	in := In{}
	e := json.Unmarshal(body, &in)
	if e != nil {
		errorLog("invalid request parameter")
		w.WriteHeader(400)
		fmt.Fprint(w, "invalid request parameter")
		return
	}

	// pipeline
	//
	//  build image
	//       |
	// run container
	//       |
	// recept stdout
	imageID := fmt.Sprintf("%d", time.Now().Unix())
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	image, errorCh1 := dockerCli.Build(ctx, imageID, in.Language, in.Version, in.Program)
	stdout, errorCh2 := dockerCli.Run(ctx, image)

	select {
	case out.Stdout = <-stdout:
		dockerCli.Clear(imageID)
	case e = <-errorCh1:
		e = dockerCli.Clear(imageID)
		out.Error = e.Error()
		if e != nil {
			out.Error = e.Error()
		}
	case e = <-errorCh2:
		out.Error = e.Error()
		dockerCli.Clear(imageID)
		if e != nil {
			out.Error = e.Error()
		}
	}
}

func log(msg string) {
	fmt.Printf("[%s] %s\n", time.Now(), msg)
}

func errorLog(msg string) {
	fmt.Fprintf(os.Stderr, "[%s] %s\n", time.Now(), msg)
}
