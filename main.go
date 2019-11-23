package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type settings struct {
	Containers map[string][]string `json:"containers"`
	Commands   map[string][]string `json:"commands"`
}

func main() {
	containerArg, commandArg := parseArgs()
	ctx := context.Background()
	cli, err := client.NewClient(client.DefaultDockerHost, client.DefaultVersion, nil, nil)
	if err != nil {
		panic(err)
	}

	s, err := read()

	containers, err := getContainerIDs(context.Background(), cli, s.Containers[containerArg])
	if err != nil {
		panic(err)
	}
	for _, container := range containers {
		if commands, ok := s.Commands[commandArg]; ok {
			for _, command := range commands {
				data, err := exec(ctx, cli, container, command)
				if err != nil {
					panic(err)
				}
				fmt.Println(string(data))
			}
		}
	}
	fmt.Println(containers)
}

func read() (*settings, error) {
	path := os.Getenv("INVOKER_SETTINGS")
	log.Print(path)
	if path == "" {
		return nil, errors.New("empty")
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var s settings
	err = json.Unmarshal(content, &s)
	return &s, err
}

func parseArgs() (string, string) {
	if len(os.Args) != 3 {
		return "", ""
	}
	return os.Args[1], os.Args[2]
}

func exec(ctx context.Context, cli *client.Client, container string, command string) ([]byte, error) {
	exec, err := cli.ContainerExecCreate(ctx, container, types.ExecConfig{
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          strings.Split(command, " "),
	})
	if err != nil {
		return nil, err
	}
	resp, err := cli.ContainerExecAttach(ctx, exec.ID, types.ExecConfig{
		AttachStderr: true,
		AttachStdout: true,
		Cmd:           strings.Split(command, " "),
	})
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Reader)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func getContainerIDs(ctx context.Context, cli *client.Client, containers []string) ([]string, error) {
	images, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, img := range images {
		for _, name := range img.Names {
			for _, container := range containers {
				if container == removePrefix(name) {
					ids = append(ids, img.ID)
				}
			}
		}
	}

	return ids, nil
}

func removePrefix(name string) string {
	return strings.Replace(name, "/", "", -1)
}
