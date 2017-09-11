package dockerCli

import ()

// BuildAndStart is
func BuildAndStart(lang, version, code string) <-chan string {
	container := make(chan string)
	go func() {
		image := build(lang, version, code)
		con := start(image)
		container <- con
		close(container)
	}()
	return container
}

// ExecuteAndStop is
func ExecuteAndStop(container <-chan string) <-chan string {
	stdout := make(chan string)

	go func() {
		con := <-container
		stdout <- execute(con)
		stop(con)
		close(stdout)
	}()

	return stdout
}

func build(lang, version, code string) {}

func start(image string) string {}

func stop(container string) int {}

func execute(container string) string {}
