package manager

// jobの管理/コンテナの管理を行う。
// apiサーバーとしても動作

// リクエストを受け取り、コンテナを選択し、jobを実行する(goroutine)(チャンネルで待つ。)
// goroutine内でコンテナを起動し、コードの実行、結果をチャンネル経由で渡す
// リクエストを返す

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
	Code     string `json:"code"`
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
	image := dockerCli.Build(in.Language, in.Version, in.Code)
	stdout := dockerCli.Run(image)
	out.Stdout = <-stdout
}
