package pipeline

import (
	"bkite/ci/pipeline/pkg/parsers"
	"fmt"
	"os"

	"github.com/buildkite/buildkite-sdk/sdk/go/sdk/buildkite"
)

type BuildContext struct {
	branch string
	commit string
}

const (
	adderPath  = "services/adder"
	subberPath = "services/adder"
)

var definedPaths = []string{adderPath, subberPath}

func ParseEvent(pipeline *buildkite.Pipeline) {
	bctx := BuildContext{
		branch: os.Getenv("BUILDKITE_BRANCH"),
		commit: os.Getenv("BUILDKITE_COMMIT"),
	}
	pull_request := os.Getenv("BUILDKITE_PULL_REQUEST")
	changeMap := parsers.ChangedPaths(definedPaths)

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
