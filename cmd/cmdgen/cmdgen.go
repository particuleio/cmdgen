package main

import (
	"fmt"
	"log"
	"os"

	"github.com/particuleio/cmdgen/pkg/cmdgen"
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
  PATH      path of the yaml config file`, argList)

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
		fmt.Printf("Missing path.\n%s\n", msgErr)
		return
	}

	path := os.Args[2]

	if err := function(path); err != nil {
		log.Fatalf(err.Error())
		return
	}

}
