package dockercli

import (
	"os/exec"
	"strconv"
)

// return image id
func Build(imageName string, dockerfilePath string) int64 {
	out, err := exec.Command("docker", "build", dockerfilePath, "-t", imageName).Output()
	if err != nil {
		return 0
	}

	imageID, err := strconv.ParseInt(string(out), 10, 64)
	if err != nil {
		return 0
	}

	return imageID
}
