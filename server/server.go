package server

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

var (
	ErrInvalidCommand = errors.New("command invalid")
)

type Response struct {
	StdOutput string `json:"std_output"`
	StdError string `json:"std_error"`
	Error string `json:"error"`
}

func run(command string) (*Response, error) {
	var resp = &Response{}
	c := strings.Split(command, " ")
	if len(c) > 1 {
		args := c[1:]
		cmd := exec.Command(c[0], args...)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			oe := fmt.Errorf("cmd.Run() failed with %s\n", err)
			resp.Error = oe.Error()
			return resp, oe
		}
		outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
		resp.StdOutput = outStr
		resp.StdError = errStr
		return resp, nil
	}
	resp.Error = ErrInvalidCommand.Error()
	return resp, ErrInvalidCommand
}