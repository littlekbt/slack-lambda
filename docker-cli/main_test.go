package dockerCli

import (
	"fmt"

	"testing"
)

func TestBuildAndRun(t *testing.T) {
	c := BuildAndRun("ruby", "2.3.0", "print 'hello world by ruby'")
	d := BuildAndRun("golang", "1.8", "package main\nimport \"fmt\"\nfunc main() {\nfmt.Println(\"hello world by golang\")\n}")
	fmt.Println(<-c)
	fmt.Println(<-d)
}
