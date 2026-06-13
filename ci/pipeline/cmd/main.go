package main

import (
	"bkite/ci/pipeline/pkg/pipeline"
	"log"
	"os"

	"github.com/buildkite/buildkite-sdk/sdk/go/sdk/buildkite"
)

func main() {
	pipe := buildkite.NewPipeline()
	pipeline.ParseEvent(pipe)
	err := finished(pipe)
	if err != nil {
		panic(err)
	}
}

func finished(pipe *buildkite.Pipeline) error {
	yaml, err := pipe.ToYAML()
	if err != nil {
		log.Fatalf("Failed to serialize YAML: %v", err)
		return err
	}

	file, err := os.Create("pipe.yaml")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(yaml)
	return err
}
