package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/buildkite/buildkite-sdk/sdk/go/sdk/buildkite"
)

type BuildContext struct {
	branch string
	commit string
}

func main() {
	pipeline := buildkite.NewPipeline()
	parse_event(pipeline)
	err := finished(pipeline)
	if err != nil {
		panic(err)
	}
}

func parse_event(pipeline *buildkite.Pipeline) {
	bctx := BuildContext{
		branch: os.Getenv("BUILDKITE_BRANCH"),
		commit: os.Getenv("BUILDKITE_COMMIT"),
	}
	pull_request := os.Getenv("BUILDKITE_PULL_REQUEST")
	changeMap := getChangedPaths()

	if pull_request != "false" {
		handlePullRequest(pipeline)
	}

	if bctx.branch == "main" {
		if changeMap[adderPath] {
			deployAdder(bctx, pipeline)
		}
		if changeMap[subberPath] {
			deploySubber(bctx, pipeline)
		}
	}
}

const (
	adderPath  = "services/adder"
	subberPath = "services/adder"
)

var definedPaths = []string{adderPath, subberPath}

func getChangedPaths() map[string]bool {
	changedPaths := strings.Split(os.Getenv("CHANGED_PATHS"), ",")
	changeMap := map[string]bool{}
	for _, changedPath := range changedPaths {
		for _, definedPath := range definedPaths {
			if strings.HasPrefix(changedPath, definedPath) {
				changeMap[definedPath] = true
			}
		}
	}
	return changeMap
}

func finished(pipe *buildkite.Pipeline) error {
	yaml, err := pipe.ToYAML()
	if err != nil {
		log.Fatalf("Failed to serialize YAML: %v", err)
		return err
	}

	file, err := os.Create("cpipeline.yaml")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(yaml)
	return err
}

func deployAdder(bctx BuildContext, pipe *buildkite.Pipeline) {
	pipe.AddStep(buildkite.CommandStep{
		Label: buildkite.Value("Deploying Adder"),
		Command: &buildkite.CommandStepCommand{
			String: buildkite.Value(fmt.Sprintf("echo 'Deploying commit %s'", bctx.commit)),
		},
	})
}

func deploySubber(bctx BuildContext, pipe *buildkite.Pipeline) {
	pipe.AddStep(buildkite.CommandStep{
		Label: buildkite.Value("Deploying Subber"),
		Command: &buildkite.CommandStepCommand{
			String: buildkite.Value(fmt.Sprintf("echo 'Deploying commit %s'", bctx.commit)),
		},
	})
}

func handlePullRequest(pipe *buildkite.Pipeline) {
	pipe.AddStep(buildkite.CommandStep{
		Label: buildkite.Value("Build"),
		Key:   buildkite.Value("build"),
		Command: &buildkite.CommandStepCommand{
			String: buildkite.Value("go build main.go"),
		},
		Plugins: &buildkite.Plugins{
			PluginsList: &buildkite.PluginsList{
				buildkite.PluginsListItem{
					String: buildkite.Value("docker#v5.13.0"),
					PluginsList: &buildkite.PluginsListObject{
						"image": "golang:1.26-alpine",
					},
				},
			},
		},
	})

	pipe.AddStep(buildkite.CommandStep{
		Label: buildkite.Value(":golang: Lint"),
		Key:   buildkite.Value("Lint"),
		Plugins: &buildkite.Plugins{
			PluginsList: &buildkite.PluginsList{
				buildkite.PluginsListItem{
					String: buildkite.Value("golangci-lint#v1.0.0"),
				},
			},
		},
	})
}
