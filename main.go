package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type settings struct {
	Containers map[string][]string `json:"containers"`
	Commands   map[string][]string `json:"commands"`
}

func main() {
	ctx := context.Background()
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	hc := &http.Client{Transport: tr}
	cli, err := client.NewClient(client.DefaultDockerHost, client.DefaultVersion, hc, nil)
	if err != nil {
		panic(err)

	}

	containers, err := getContainerIDs(context.Background(), cli, []string{"api_couchbase_1"})
	if err != nil {
		panic(err)
	}
	for _, container := range containers {
		data, err:= exec(ctx, cli, container)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(data))
	}
	fmt.Println(containers)
}

func exec(ctx context.Context, cli *client.Client, container string) ([]byte, error) {
	exec, err := cli.ContainerExecCreate(ctx, container, types.ExecConfig{
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          []string{"ls", "-ll"},
	})
	if err != nil {
		return nil, err
	}
	resp, err := cli.ContainerExecAttach(ctx, exec.ID, types.ExecConfig{
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          []string{"ls", "-ll"},
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
