package main

import (
	"fmt"
	"log"
	"os"

	"github.com/buildkite/buildkite-sdk/sdk/go/sdk/buildkite"
)

type BuildContext struct {
	branch string
	commit string
}

func main() {
	buildContext := BuildContext{
		branch: os.Getenv("BUILDKITE_BRANCH"),
		commit: os.Getenv("BUILDKITE_COMMIT"),
	}
	log.Printf("BUILDKITE_GITHUB_ACTION: %s", os.Getenv("BUILDKITE_GITHUB_ACTION"))
	log.Printf("BUILDKITE_GITHUB_EVENT: %s", os.Getenv("BUILDKITE_GITHUB_EVENT"))
	log.Printf("BUILDKITE_MESSAGE: %s", os.Getenv("BUILDKITE_MESSAGE"))
	log.Printf("BUILDKITE_PULL_REQUEST: %s", os.Getenv("BUILDKITE_PULL_REQUEST"))
	pipeline := buildkite.Pipeline{}

	githubEvent := os.Getenv("BUILDKITE_GITHUB_EVENT")
	if githubEvent == "pull_request" {
		pipeline = handlePullRequest(pipeline)
	} else if buildContext.branch == "main" {
		pipeline = handlePush(buildContext, pipeline)
	} else {
		log.Fatal("Unknown event")
	}

	yaml, err := pipeline.ToYAML()
	if err != nil {
		log.Fatalf("Failed to serialize YAML: %v", err)
	}

	fmt.Println(yaml)
}

func handlePush(bctx BuildContext, pipe buildkite.Pipeline) buildkite.Pipeline {
	pipe.AddStep(buildkite.CommandStep{
		Command: &buildkite.CommandStepCommand{
			String: buildkite.Value(fmt.Sprintf("echo 'Deploying commit %s'", bctx.commit)),
		},
	})

	return pipe
}

func handlePullRequest(pipe buildkite.Pipeline) buildkite.Pipeline {
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
	return pipe
}
