package main

import (
	"fmt"
	"log"

	"github.com/buildkite/buildkite-sdk/sdk/go/sdk/buildkite"
)

func main() {
	pipeline := buildkite.Pipeline{}

	pipeline.AddStep(buildkite.CommandStep{
		Command: &buildkite.CommandStepCommand{
			String: buildkite.Value("echo 'Hello, world!"),
		},
	})

	// YAML output
	yaml, err := pipeline.ToYAML()
	if err != nil {
		log.Fatalf("Failed to serialize YAML: %v", err)
	}

	fmt.Println(yaml)
}
