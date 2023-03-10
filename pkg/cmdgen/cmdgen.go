package cmdgen

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"

	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

type cmdItem struct {
	Cmd         string `yaml:"cmd"`
	Description string `yaml:"description"`
}

type templateStructure struct {
	Scenario []cmdItem `yaml:"scenario"`
	Clean    []string  `yaml:"clean"`
}

const (
	ScriptFileExt = ".bash"
	ShellToUse    = "bash"
	Template      = `#!/bin/bash

{{ range .Scenario }}
	{{- range (split .Description) -}}
		{{- printf "# %s\n" . }}
	{{- end -}}
	{{- println .Cmd }}
{{ end }} 

## To clean workspace run these commands:
{{- range .Clean }}
# {{ . }}
{{- end }}

`
)

var (
	ErrFoundDir      = errors.New("found directory expected file")
	colorHeadline    = color.New(color.FgHiBlack)
	colorDescription = color.New(color.FgGreen)
	colorCmd         = color.New(color.FgYellow)
	colorOutput      = color.New(color.FgCyan)
	colorError       = color.New(color.FgRed)
)

func (c cmdItem) String() string {
	return fmt.Sprintf(`{"cmd": "%s", "description": "%s"}`, c.Cmd, c.Description)
}

func (c cmdItem) printDescription() {
	colorHeadline.Println("[description] ")
	colorDescription.Println(c.Description)
}

func (c cmdItem) printCmd() {
	colorHeadline.Println("[cmd] ")
	colorCmd.Println(c.Cmd)
}

func checkFile(filePath string) error {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return err
	}

	if info.IsDir() {
		return ErrFoundDir
	}

	return nil
}

func parseFile(filePath string) (ts templateStructure, err error) {
	if err = checkFile(filePath); err != nil {
		return
	}

	var content []byte
	if content, err = os.ReadFile(filePath); err != nil {
		return
	}

	err = yaml.Unmarshal(content, &ts)
	return
}

func createFile(filePath string) (file *os.File, err error) {
	fileName := strings.Split(path.Base(filePath), ".")[0] + ScriptFileExt
	if file, err = os.Create(fileName); err != nil {
		return
	}

	return
}

func writeFormatted(file io.Writer, ts templateStructure) (err error) {
	t, err := template.New("script").Funcs(template.FuncMap{
		"split": func(text string) []string {
			return strings.Split(text, "\n")
		},
	}).Parse(Template)
	if err != nil {
		return
	}
	wr := bufio.NewWriter(file)
	if err = t.Execute(wr, ts); err != nil {
		return
	}
	wr.Flush()
	return
}

func GenBashScript(filePath string) (err error) {
	// parse template file
	var ts templateStructure
	if ts, err = parseFile(filePath); err != nil {
		return
	}

	// create new file
	var file *os.File
	if file, err = createFile(filePath); err != nil {
		return
	}
	defer file.Close()

	// Write formatted template
	writeFormatted(file, ts)
	return nil
}

func printStd(std io.ReadCloser, c *color.Color) {
	outputScanner := bufio.NewScanner(std)
	outputScanner.Split(bufio.ScanLines)
	for outputScanner.Scan() {
		output := outputScanner.Text()
		c.Println(output)
	}
}

func printOutput(stdout, stderr io.ReadCloser) {
	colorHeadline.Println("[output] ")
	printStd(stdout, colorOutput)
	printStd(stderr, colorError)
}

func printItem(index int, cmd cmdItem, stdout, stderr io.ReadCloser) {
	colorHeadline.Printf("[Step %d]:", index)
	fmt.Println()
	cmd.printDescription()
	fmt.Println()
	if len(strings.TrimSpace(cmd.Cmd)) > 0 {
		cmd.printCmd()
		fmt.Println()
		printOutput(stdout, stderr)
	}
}

func processCmd(index int, cmd cmdItem) error {
	c := exec.Command(ShellToUse, "-c", cmd.Cmd)
	stdout, _ := c.StdoutPipe()
	stderr, _ := c.StderrPipe()
	c.Start()

	printItem(index, cmd, stdout, stderr)

	c.Wait()

	return nil
}

func CleanWorkspace(filePath string) (err error) {
	ts, err := parseFile(filePath)
	for _, cmd := range ts.Clean {
		exec.Command(ShellToUse, "-c", cmd).Run()
	}
	return
}

func StartScenario(filePath string) (err error) {
	ts, err := parseFile(filePath)

	for index, cmd := range ts.Scenario {
		fmt.Println("-------------")
		processCmd(index+1, cmd)
		fmt.Println()
		fmt.Print("press <Enter> to conitnue ")
		fmt.Scanln()
	}

	fmt.Println("END OF PROCESS")
	return
}
