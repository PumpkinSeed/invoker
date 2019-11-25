package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

const (
	envPrefix = "INVOKER_"

	envVerbose    = envPrefix + "VERBOSE"
	envSkipOutput = envPrefix + "SKIP_OUTPUT"
	envSettings   = envPrefix + "SETTINGS"
)

type settings struct {
	Containers map[string][]string `json:"containers"`
	Commands   map[string][]string `json:"commands"`
}

func main() {
	containerArg, commandArg := parseArgs()
	verboseLog(fmt.Sprintf("invoker starts to execute {%s} command group on {%s} container group", commandArg, containerArg), info)

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		verboseLog(fmt.Sprintf("%s [%s]", "error: creating Docker client connection", err.Error()), fatal)
		return
	}

	s := read()

	containers, err := getContainerIDs(context.Background(), cli, s.Containers[containerArg])
	if err != nil {
		panic(err)
	}
	for _, container := range containers {
		if commands, ok := s.Commands[commandArg]; ok {
			for _, command := range commands {
				data := exec(ctx, cli, container.id, command)
				printOutput(container, command, string(data))
			}
		}
	}
}

func printOutput(container containerDef, command, data string) {
	if skip, ok := os.LookupEnv(envSkipOutput); ok && skip == "true" {
		return
	}
	fmt.Printf("%s [%s] ~ %s\n", removePrefix(container.name), container.id, command)
	fmt.Print(data)
}

func read() *settings {
	path := os.Getenv(envSettings)
	if path == "" {
		verboseLog(fmt.Sprintf("%s", "error: settings path is empty, set INVOKER_SETTINGS with the path of settings file"), fatal)
		return nil
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		verboseLog(fmt.Sprintf("%s %s [%s]", "error: reading the file on path", path, err.Error()), fatal)
		return nil
	}

	var s settings
	err = json.Unmarshal(content, &s)
	if err != nil {
		verboseLog(fmt.Sprintf("%s [%s]", "error: marshal the file content to settings", err.Error()), fatal)
		return nil
	}
	return &s
}

func parseArgs() (string, string) {
	if len(os.Args) != 3 {
		return "", ""
	}
	return os.Args[1], os.Args[2]
}

func exec(ctx context.Context, cli *client.Client, container string, command string) []byte {
	exec, err := cli.ContainerExecCreate(ctx, container, types.ExecConfig{
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          strings.Split(command, " "),
	})
	if err != nil {
		verboseLog(fmt.Sprintf("%s [%s]", "error: create exec in container", err.Error()), fatal)
		return nil
	}

	resp, err := cli.ContainerExecAttach(ctx, exec.ID, types.ExecConfig{
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          strings.Split(command, " "),
	})
	if err != nil {
		verboseLog(fmt.Sprintf("%s [%s]", "error: attach exec to container", err.Error()), fatal)
		return nil
	}
	data, err := ioutil.ReadAll(resp.Reader)
	if err != nil {
		verboseLog(fmt.Sprintf("%s [%s]", "error: parse the response", err.Error()), fatal)
		return nil
	}

	return data
}

type containerDef struct {
	id   string
	name string
}

func getContainerIDs(ctx context.Context, cli *client.Client, containers []string) ([]containerDef, error) {
	images, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	var ids []containerDef
	for _, img := range images {
		for _, name := range img.Names {
			for _, container := range containers {
				if container == removePrefix(name) {
					ids = append(ids, containerDef{
						id:   img.ID,
						name: name,
					})
				}
			}
		}
	}

	return ids, nil
}

func removePrefix(name string) string {
	return strings.Replace(name, "/", "", -1)
}

const (
	info = iota
	fatal
)

func verboseLog(str string, level int) {
	if verbose, ok := os.LookupEnv(envVerbose); ok && verbose != "true" {
		return
	} else if !ok {
		return
	}

	switch level {
	case info:
		log.Printf("[%s] %s", "INFO", str)
	case fatal:
		log.Printf("[%s] %s", "FATAL", str)
		os.Exit(1)
	}
}
