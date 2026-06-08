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
	changes := os.Getenv("CHANGED_PATHS")
	changeMap := parsers.ChangedPaths(changes, definedPaths)

	fmt.Println(pull_request)
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
	for _, v := range []string{"services/adder/cmd/main.go", "services/subber/cmd/main.go"} {
		buildService(pipe, v)
	}

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

	pipe.AddStep(buildkite.CommandStep{
		Label: buildkite.Value("Test"),
		Key:   buildkite.Value("Test"),
		Command: &buildkite.CommandStepCommand{
			String: buildkite.Value("./test.sh gotestsum --junitfile junit.xml ./apitest/..."),
		},
		Secrets: &buildkite.Secrets{
			Secrets: &buildkite.SecretsObject{
				"BUILDKITE_ANALYTICS_TOKEN": "TEST_SUITE_KEY",
			},
		},
		Plugins: &buildkite.Plugins{
			PluginsList: &buildkite.PluginsList{
				buildkite.PluginsListItem{
					PluginsList: &buildkite.PluginsListObject{
						"test-collector#v1.11.0": map[string]string{
							"format": "junit",
							"files":  "junit.xml",
						},
					},
				},
			},
		},
	})
}

func buildService(pipe *buildkite.Pipeline, service string) {
	pipe.AddStep(buildkite.CommandStep{
		Label: buildkite.Value("Build"),
		Command: &buildkite.CommandStepCommand{
			String: buildkite.Value(fmt.Sprintf("go build %s", service)),
		},
		Plugins: &buildkite.Plugins{
			PluginsList: &buildkite.PluginsList{
				{
					PluginsList: &buildkite.PluginsListObject{
						"docker#v5.13.0": map[string]string{
							"image": "golang:1.26-alpine",
						},
					},
				},
			},
		},
	})
}
