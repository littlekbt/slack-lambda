package dockercli

import (
	"testing"
)

var expectJsDockerFile = `
FROM node:6.10 
MAINTAINER littlekbt

WORKDIR src/

COPY /tmp/1/main.js main.js

CMD ["node", "main.js"]
`

func TestMakeDockerFile(t *testing.T) {
	c := langToConfig["node"]
	actual := c.compile("/tmp/1/main.js")
	if actual != expectJsDockerFile {
		t.Errorf("actual: %s \n expected: %s", actual, expectJsDockerFile)
	}
}
