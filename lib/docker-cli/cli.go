package dockercli

import (
	"fmt"
	"os/exec"
)

// この辺の情報はmanager側に移動
var langToExtension = map[string]string{
	"node": "js",
}

var langToImage = map[string]string{
	"node": "node6.10",
}

// idの発行は外部で(mysqlに依存させたくない)
// dockerfileも外で
// dockerfileのパスとimageName受け取ってbuildするだけ

// return image id
func Build(imageName string, dockerfilePath string) int {
	out, err := exec.Command("docker", "build", dockerfilePath, "-t", imageName).Output()
	if err != nil {
		return 0
	}

	imageID, err := strconv.ParseInt(out, 10, 64)
	if err != nil {
		return 0
	}

	return imageID
}

// 実行のみ

// return stdout
func Run(imageId int) string {}

func Rm(containerID string) bool {}

func Rmi(imageID string) bool {}
