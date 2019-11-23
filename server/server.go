package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

const (
	DefaultPort = 3000
)

var (
	ErrInvalidCommand = errors.New("command invalid")
)

type Response struct {
	CommandResponses []CommandResponse `json:"command_responses"`
	Error string `json:"error"`
}

type CommandResponse struct {
	Iteration int `json:"iteration"`
	Name string `json:"name"`
	Command string `json:"command"`
	StdOutput string `json:"std_output"`
	StdError string `json:"std_error"`
	Error string `json:"error"`
}

type settings struct {
	Port int `json:"port"`
	Commands map[string][]string `json:"commands"`
}

type server struct {
	path string
}

func Serve(s string) {
	server := server{path: s}
	http.HandleFunc("/", server.handler)
	log.Printf("Listening on :%d", DefaultPort)
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}

func (s *server) handler(w http.ResponseWriter, r *http.Request) {
	commands := strings.Split(r.URL.Path[1:], "/")
	log.Print(commands)
	settings, err := readSettings(s.path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(Response{Error: fmt.Sprintf("Error while reading settings [%s]", err.Error())})
		_, _ = w.Write(resp)
		return
	}

	var resp Response
	for _, command := range commands {
		if v, ok := settings.Commands[command]; ok {
			for i, actualCommand := range v {
				cr := run(actualCommand)
				cr.Iteration = i
				cr.Name = command
				resp.CommandResponses = append(resp.CommandResponses, cr)
			}
		}
	}
	w.WriteHeader(http.StatusOK)
	jresp, _ := json.Marshal(resp)
	_, _ = w.Write(jresp)
	return
}

func readSettings(s string) (*settings, error) {
	dat, err := ioutil.ReadFile(s)
	if err != nil {
		return nil, err
	}
	var settings settings
	if err := json.Unmarshal(dat, &settings); err != nil {
		return nil, err
	}

	return &settings, nil
}

func run(command string) CommandResponse {
	var resp = CommandResponse{
		Command: command,
	}
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
			return resp
		}
		outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
		resp.StdOutput = outStr
		resp.StdError = errStr
		return resp
	}
	resp.Error = ErrInvalidCommand.Error()
	return resp
}