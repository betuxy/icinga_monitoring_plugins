package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/rs/zerolog"
)

const (
	ICINGA_OK       int = 0
	ICINGA_CRITICAL int = 2
	ICINGA_UNKNOWN  int = 3
)

func interpretResults(containerList []map[string]string, containerExcludes []string) {
	containersNotRunning := []map[string]string{}

	for _, container := range containerList {
		if container["state"] != "running" {
			//fmt.Println(containerExcludes, container["name"])
			if !slices.ContainsFunc(containerExcludes, func(s string) bool {
				return strings.Contains(container["name"], s)
			}) {
				containersNotRunning = append(containersNotRunning, map[string]string{
					container["name"]: container["state"],
				})
			}
		}
	}

	if len(containersNotRunning) > 0 {
		fmt.Printf("[ERROR] The following containers are not running:\n")
		for _, container := range containersNotRunning {
			for name, state := range container {
				fmt.Printf("  - %v: %s\n", name, state)
			}
		}
		os.Exit(ICINGA_CRITICAL)

	} else {
		fmt.Printf("[OK] All containers are running.")
		os.Exit(ICINGA_OK)
	}
}

func main() {
	deploymentPrefix := flag.String("p", "", "Container name prefix (generally applicable in compose deployments)")
	debugLog := flag.Bool("d", false, "Debug verbosity")

	var containerNames containerNamesList
	flag.Var(&containerNames, "n", "Specify name or substring of container to add to deployment (can be set multiple times)")

	var containerExcludes containerExcludesList
	flag.Var(&containerExcludes, "x", "Specify name or substring of container to exclude from check (can be set multiple times)")

	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if *debugLog {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if len(*deploymentPrefix) > 0 && len(containerNames) > 0 {
		var allContainers []map[string]string
		containers, err := getContainersByPrefix(*deploymentPrefix)
		if err != nil {
			fmt.Println("[ERROR]", err)
			os.Exit(ICINGA_UNKNOWN)
		}

		allContainers = slices.Concat(allContainers, containers)

		containers, err = getContainersByName(containerNames)
		if err != nil {
			fmt.Println("[ERROR]", err)
			os.Exit(ICINGA_UNKNOWN)
		}

		allContainers = slices.Concat(allContainers, containers)
		interpretResults(allContainers, containerExcludes)

	} else if len(*deploymentPrefix) > 0 && len(containerNames) < 1 {
		containers, err := getContainersByPrefix(*deploymentPrefix)
		if err != nil {
			fmt.Println("[ERROR]", err)
			os.Exit(ICINGA_UNKNOWN)
		}

		interpretResults(containers, containerExcludes)

	} else if len(*deploymentPrefix) < 1 && len(containerNames) > 0 {
		containers, err := getContainersByName(containerNames)
		if err != nil {
			fmt.Println("[ERROR]", err)
			os.Exit(ICINGA_UNKNOWN)
		}

		interpretResults(containers, containerExcludes)
	}
}
