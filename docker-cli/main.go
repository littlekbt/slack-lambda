package dockerCli

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"
)

var langToCommand = map[string]string{
	"ruby":   "\"ruby\"",
	"Ruby":   "\"ruby\"",
	"go":     "\"go\",\"run\"",
	"golang": "\"go\",\"run\"",
}

var langToExt = map[string]string{
	"ruby":   "rb",
	"golang": "go",
}

var langToImageName = map[string]string{
	"ruby":   "ruby",
	"Ruby":   "ruby",
	"go":     "golang",
	"golang": "golang",
}

func init() {
	_, err := exec.LookPath("docker")

	if err != nil {
		panic("not found docker command")
	}
}

// Build image
func Build(lang, version, code string) <-chan string {
	container := make(chan string)
	go func() {
		uuid := fmt.Sprintf("%d_%s", time.Now().Unix(), code[0:5])

		dir, err := mkDirP(uuid)
		if err != nil {
			panic("fail mkdir")
		}

		programfile, err := createTmpProgramFile(dir, lang, code)
		if err != nil {
			panic("fail create program file")
		}

		dockerfile, err := createTmpDockerFile(dir, lang, version, programfile)
		if err != nil {
			panic("fail create Dockerfile")
		}

		image, err := build(path.Dir(dockerfile), uuid)
		if err != nil {
			panic("fail build image")
		}

		container <- image
		close(container)
	}()

	return container
}

// Run executes container
func Run(image <-chan string) chan string {
	stdout := make(chan string)
	go func() {
		imageName := <-image
		out, err := run(imageName)
		if err != nil {
			panic("fail run container")
		}
		stdout <- out

		// side effects...
		clear(imageName)
	}()

	return stdout
}

// clear removes container, image, files
func clear(image string) {
	err1 := exec.Command("sh", "-c", fmt.Sprintf("docker rm -f `docker ps -a | grep %s | awk '{print $1}'`", image)).Run()
	if err1 != nil {
		panic(err1.Error())
	}
	err2 := exec.Command("sh", "-c", fmt.Sprintf("docker rmi -f `docker images | grep %s | awk '{print $3}'`", image)).Run()
	if err2 != nil {
		panic(err2.Error())
	}
	err3 := os.RemoveAll(path.Join("/tmp", image))
	if err3 != nil {
		panic(err3.Error())
	}
}

// return directory
func mkDirP(uuid string) (string, error) {
	dir := fmt.Sprintf("/tmp/slack-lambda/%s", uuid)
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return "", err
	}

	return dir, nil
}

// return programfile path
func createTmpProgramFile(dir, lang, code string) (string, error) {
	ext := langToExt[lang]
	filepath := path.Join(dir, fmt.Sprintf("program.%s", ext))

	file, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	file.Write([]byte(code))

	return filepath, nil
}

// return dockerfile path
func createTmpDockerFile(dir, lang, version, programfilePath string) (string, error) {
	dockerfilepath := path.Join(dir, "Dockerfile")

	file, err := os.Create(dockerfilepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	langCmd := langToCommand[lang]
	ext := langToExt[lang]
	image := langToImageName[lang]
	baseimage := fmt.Sprintf("%s:%s", image, version)

	dockerfile := fmt.Sprintf(`
FROM %s

WORKDIR /src

COPY . .

CMD [%s, "program.%s"]
  `, baseimage, langCmd, ext)

	file.Write([]byte(dockerfile))

	return dockerfilepath, nil
}

func build(dir, uuid string) (string, error) {
	image := fmt.Sprintf("slack-lambda/%s", uuid)
	err := exec.Command("docker", "build", dir, "-t", image).Run()
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return image, nil
}

func run(image string) (string, error) {
	out, err := exec.Command("docker", "run", image).Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
