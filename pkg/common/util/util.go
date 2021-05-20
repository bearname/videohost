package util

import (
	"bytes"
	"fmt"
	"os/exec"
)

func CmdExec(args ...string) (string, error) {
	baseCmd := args[0]
	cmdArgs := args[1:]

	cmd := exec.Command(baseCmd, cmdArgs...)
	var outputBuffer, errorBuffer bytes.Buffer
	cmd.Stdout = &outputBuffer
	cmd.Stderr = &errorBuffer
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println("out:", outputBuffer.String(), "err:", errorBuffer.String())

	return outputBuffer.String(), nil
}
