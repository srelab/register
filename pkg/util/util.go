package util

import (
	"bytes"
	"os/exec"
	"strings"
)

func StripIndent(multilineStr string) string {
	return strings.Replace(multilineStr, "\t", "", -1)
}

func CmdOutBytes(name string, arg ...string) ([]byte, error) {
	cmd := exec.Command(name, arg...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.Bytes(), err
}
