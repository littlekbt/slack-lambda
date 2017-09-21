package dockerCli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
)

var langToCommand = map[string]string{
	"ruby":   `"ruby"`,
	"Ruby":   `"ruby"`,
	"go":     `"go", "run"`,
	"golang": `"go", "run"`,
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
func Build(ctx context.Context, imageID, lang, version, program string) (<-chan string, chan error) {
	container := make(chan string)
	errCh := make(chan error)

	go func() {
		dir, err := mkDirP(imageID)
		if err != nil {
			errCh <- errors.New("fail mkdir")
		}

		programfile, err := createTmpProgramFile(dir, lang, program)
		if err != nil {
			errCh <- errors.New("fail create program file")
		}

		dockerfile, err := createTmpDockerFile(dir, lang, version, programfile)
		if err != nil {
			errCh <- errors.New("fail create Dockerfile")
		}

		image, err := build(path.Dir(dockerfile), imageID)
		if err != nil {
			errCh <- errors.New("fail build image. stderr message: " + err.Error())
		}

		select {
		case <-ctx.Done():
			errCh <- errors.New("timeout")
			return
		case container <- image:
			close(container)
		}
	}()

	return container, errCh
}

// Run executes container
func Run(ctx context.Context, image <-chan string) (<-chan string, chan error) {
	stdout := make(chan string)
	errCh := make(chan error)

	go func() {
		imageName := <-image
		out, err := run(imageName)
		if err != nil {
			errCh <- errors.New("fail run container.\n Error message: " + out)
		}

		select {
		case <-ctx.Done():
			errCh <- errors.New("timeout")
			return
		case stdout <- out:
			close(stdout)
		}
	}()

	return stdout, errCh
}

// Clear removes container, image, files
func Clear(imageID string) error {
	err1 := exec.Command("sh", "-c", fmt.Sprintf("docker rm -f `docker ps -a | grep %s | awk '{print $1}'`", imageID)).Run()
	if err1 != nil {
		return err1
	}
	err2 := exec.Command("sh", "-c", fmt.Sprintf("docker rmi -f `docker images | grep %s | awk '{print $3}'`", imageID)).Run()
	if err2 != nil {
		return err2
	}
	err3 := os.RemoveAll(path.Join("/tmp", imageID))
	if err3 != nil {
		return err3
	}

	return nil
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
func createTmpProgramFile(dir, lang, program string) (string, error) {
	ext := langToExt[lang]
	filepath := path.Join(dir, fmt.Sprintf("program.%s", ext))

	file, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	file.Write([]byte(program))

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
	out, err := exec.Command("sh", "-c", fmt.Sprintf("docker run %s 2>&1", image)).Output()
	if err != nil {
		return string(out), err
	}

	return string(out), nil
}
