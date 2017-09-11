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
	Code     string `json:"Code"`
}

// Out is response body
type Out struct {
	Stdout string `json:"Stdout"`
	Error  string `json:"Error"`
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
		out.Error = "request is only post"
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

	// inputの内容でgoroutineを作る
	// pipeで繋ぐ
	//    コンテナ起動
	//         |
	//        実行
	//   |            |
	// 結果受け取り コンテナ停止

	ch := make(<-chan string)

	func() {
		container := dockerCli.BuildAndStart(in.Language, in.Version, in.Code)
		ch = dockerCli.ExecuteAndStop(container)
	}()

	out.Stdout = <-ch
}
