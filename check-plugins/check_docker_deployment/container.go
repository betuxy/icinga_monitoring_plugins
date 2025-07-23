package main

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/rs/zerolog/log"
)

// Variable Type for Flag containerNames
type containerNamesList []string

func (i *containerNamesList) String() string {
	stringRepresentation := ""
	for _, name := range *i {
		stringRepresentation += (name + "\n")
	}

	return stringRepresentation
}

func (i *containerNamesList) Set(value string) error {
	*i = append(*i, strings.TrimSpace(value))
	return nil
}

// Variable Type for Flag containerExcludes
type containerExcludesList []string

func (i *containerExcludesList) String() string {
	stringRepresentation := ""
	for _, name := range *i {
		stringRepresentation += (name + "\n")
	}

	return stringRepresentation
}

func (i *containerExcludesList) Set(value string) error {
	*i = append(*i, strings.TrimSpace(value))
	return nil
}

func printContainerInfo(containerList []map[string]string) {
	fmt.Printf("%-30s%-20s%s\n", "Name", "State", "ID")
	for _, container := range containerList {
		fmt.Printf("%-30v%-20v%v\n", container["name"], container["state"], container["id"])
	}
}

func getContainersByPrefix(prefix string) ([]map[string]string, error) {
	log.Debug().Msg("getContainersByPrefix:param:prefix: " + prefix)
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println("[ERROR]", err)
		os.Exit(ICINGA_UNKNOWN)
	}

	defer cli.Close()

	containerList, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		fmt.Println("[ERROR]", err)
		os.Exit(ICINGA_UNKNOWN)
	}

	containerNames := []map[string]string{}
	for _, container := range containerList {
		log.Debug().Msg("getContainersByPrefix:containerName: " + container.Names[0])
		if strings.HasPrefix(strings.TrimPrefix(container.Names[0], "/"), prefix) {
			containerNames = append(containerNames, map[string]string{
				"id":    container.ID,
				"name":  strings.TrimPrefix(container.Names[0], "/"),
				"state": container.State,
			})
		}
	}

	return containerNames, nil
}

func getContainersByName(containerNames []string) ([]map[string]string, error) {
	log.Debug().Msg("getContainersByName:")
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println("[ERROR]", err)
		os.Exit(ICINGA_UNKNOWN)
	}

	defer cli.Close()

	containerList, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		fmt.Println("[ERROR]", err)
		os.Exit(ICINGA_UNKNOWN)
	}

	containerResults := []map[string]string{}
	for _, container := range containerList {
		containerName := strings.TrimPrefix(container.Names[0], "/")
		if slices.ContainsFunc(containerNames, func(s string) bool {
			return strings.Contains(containerName, s)
		}) {

			containerResults = append(containerResults, map[string]string{
				"id":    container.ID,
				"name":  containerName,
				"state": container.State,
			})
		}
	}

	return containerResults, nil
}
