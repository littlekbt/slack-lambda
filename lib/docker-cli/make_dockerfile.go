package dockercli

import (
	"errors"
	"regexp"
)

var dockerFileTemplate = `
FROM <<baseImageName>> 
MAINTAINER littlekbt

WORKDIR src/

COPY <<programFilePath>> main.<<extention>>

CMD [<<cmd>>, "main.<<extention>>"]
`

type config struct {
	baseImageName string
	extension     string
	cmd           string
}

var langToConfig = map[string]config{
	"node": config{
		baseImageName: "node:6.10",
		extension:     "js",
		cmd:           "\"node\"",
	},
}

func (c *config) compile(programFilePath string) string {
	temp := dockerFileTemplate
	temp = replaceString(temp, "baseImageName", c.baseImageName)
	temp = replaceString(temp, "extention", c.extension)
	temp = replaceString(temp, "cmd", c.cmd)
	temp = replaceString(temp, "programFilePath", programFilePath)
	return temp
}

func replaceString(source string, cond string, changeStr string) string {
	r := regexp.MustCompile("<<" + cond + ">>")
	return r.ReplaceAllString(source, changeStr)
}

// programファイルとdockerfileの置き場を決めて渡す
// MakeDockerFile make dockerfile contents.
func MakeDockerFile(lang string, programFilePath string, dockerFilePath string) (string, error) {
	c, ok := langToConfig[lang]

	if !ok {
		return "", errors.New("no support lang")
	}

	return c.compile(programFilePath), nil
}
