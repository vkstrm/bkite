package pipeline

import (
	"bkite/ci/pipeline/pkg/parsers"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/buildkite/buildkite-sdk/sdk/go/sdk/buildkite"
)

type buildContext struct {
	branch string
	commit string
}

const (
	adderPath  = "services/adder"
	subberPath = "services/subber"
)

var definedPaths = []string{adderPath, subberPath}

func ParseEvent(pipeline *buildkite.Pipeline) {
	bctx := buildContext{
		branch: os.Getenv("BUILDKITE_BRANCH"),
		commit: os.Getenv("BUILDKITE_COMMIT"),
	}
	pull_request := os.Getenv("BUILDKITE_PULL_REQUEST")
	changes := os.Getenv("CHANGED_PATHS")
	changeMap := parsers.ChangedPaths(changes, definedPaths)

	log.Printf("BUILDKITE_PULL_REQUEST: %s", pull_request)
	_, err := strconv.Atoi(pull_request)
	if err == nil {
		handlePullRequest(pipeline)
		return
	}
	if bctx.branch == "main" {
		handleRelease(bctx, pipeline, changeMap)
	}
}

func handleRelease(bctx buildContext, pipe *buildkite.Pipeline, changeMap map[string]bool) {
	if len(changeMap) == 0 {
		log.Print("No changed paths")
		return
	}

	testStageDepends := []buildkite.DependsOnListItem{}
	if changeMap[adderPath] {
		deployAdder(bctx, pipe, "stage")
		testStageDepends = append(testStageDepends, buildkite.DependsOnListItem{
			String: buildkite.Value("adder-stage"),
		})
	}
	if changeMap[subberPath] {
		deploySubber(bctx, pipe, "stage")
		testStageDepends = append(testStageDepends, buildkite.DependsOnListItem{
			String: buildkite.Value("subber-stage"),
		})
	}
	testStep(pipe, "stage", &testStageDepends)
	prodButton(pipe)
	testProdDepends := []buildkite.DependsOnListItem{}
	if changeMap[adderPath] {
		deployAdder(bctx, pipe, "prod")
		testProdDepends = append(testProdDepends, buildkite.DependsOnListItem{
			String: buildkite.Value("adder-prod"),
		})
	}
	if changeMap[subberPath] {
		deploySubber(bctx, pipe, "prod")
		testProdDepends = append(testProdDepends, buildkite.DependsOnListItem{
			String: buildkite.Value("subber-stage"),
		})
	}
	testStep(pipe, "prod", &testProdDepends)

}

func prodButton(pipe *buildkite.Pipeline) {
	pipe.AddStep(buildkite.BlockStep{
		Block:  buildkite.Value("Ship it?"),
		Prompt: buildkite.Value("Ship it?"),
	})
}

func deployAdder(bctx buildContext, pipe *buildkite.Pipeline, environment string) {
	pipe.AddStep(buildkite.CommandStep{
		Label: buildkite.Value("Deploying Adder"),
		Key:   buildkite.Value(fmt.Sprintf("adder-%s", environment)),
		Command: &buildkite.CommandStepCommand{
			String: buildkite.Value(fmt.Sprintf("echo 'Deploying commit %s to %s'", bctx.commit, environment)),
		},
	})
}

func deploySubber(bctx buildContext, pipe *buildkite.Pipeline, environment string) {
	pipe.AddStep(buildkite.CommandStep{
		Label: buildkite.Value("Deploying Subber"),
		Key:   buildkite.Value(fmt.Sprintf("subber-%s", environment)),
		Command: &buildkite.CommandStepCommand{
			String: buildkite.Value(fmt.Sprintf("echo 'Deploying commit %s to %s'", bctx.commit, environment)),
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

	testStep(pipe, "pull-request", nil)
}

func testStep(pipe *buildkite.Pipeline, environment string, dependsOn *buildkite.DependsOnList) {
	pipe.AddStep(buildkite.CommandStep{
		Label: buildkite.Value("Test"),
		Key:   buildkite.Value(fmt.Sprintf("test-%s", environment)),
		DependsOn: &buildkite.DependsOn{
			DependsOnList: dependsOn,
		},
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
