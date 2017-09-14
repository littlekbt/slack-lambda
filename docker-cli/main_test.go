package dockerCli

import (
	"fmt"
	"testing"
)

func TestBuildAndRun(t *testing.T) {
	image1 := Build("ruby", "2.3.0", "print 'hello world by ruby'")
	image2 := Build("golang", "1.8", "package main\nimport \"fmt\"\nfunc main() {\nfmt.Println(\"hello world by golang\")\n}")
	c := Run(image1)
	d := Run(image2)
	go fmt.Println(<-c)
	go fmt.Println(<-d)
}
