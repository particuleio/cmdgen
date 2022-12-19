package main

import "cmdgen/cmd/cmdgen"

const (
	dummyTestFile = "files/test.yaml"
)

func main() {
	cmdgen.StartScenario(dummyTestFile)
}
