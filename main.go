package main

import "cmdgen/pkg/cmdgen"

const (
	dummyTestFile = "files/test.yaml"
)

func main() {
	cmdgen.StartScenario(dummyTestFile)
}
