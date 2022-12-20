package main

import (
	"cmdgen/pkg/cmdgen"
	"fmt"
	"log"
	"os"
)

func main() {
	argMap := map[string]func(string) error{
		"start":    cmdgen.StartScenario,
		"clean":    cmdgen.CleanWorkspace,
		"generate": cmdgen.GenBashScript,
	}

	argList := []string{}
	for arg := range argMap {
		argList = append(argList, arg)
	}

	msgErr := fmt.Sprintf(`Usage: cmdgen COMMAND PATH
  COMMAND   one of %v
  PATH      of the yaml config file`, argList)

	if len(os.Args) < 2 {
		fmt.Println(msgErr)
		return
	}

	function, ok := argMap[os.Args[1]]
	if !ok {
		fmt.Printf("Wrong argument.\n%s\n", msgErr)
		return
	}

	if len(os.Args) < 3 {
		fmt.Printf("missing path.\n%s\n", msgErr)
		return
	}

	path := os.Args[2]

	if err := function(path); err != nil {
		log.Fatalf(err.Error())
		return
	}

}
