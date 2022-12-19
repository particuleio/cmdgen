package cmdgen

import (
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

const (
	dummyTestFile = "../../files/test.yaml"
)

func TestTemplateFile(t *testing.T) {
	// It should exist
	info, err := os.Stat(TemplateFile)
	if os.IsNotExist(err) {
		t.Fatalf("TemplateFile does not exist")
	}

	// It should not be a directory
	if info.IsDir() {
		t.Fatalf("TemplateFile expected to be a file. Got a dir")
	}

	// It should end with TemplateFileExt
	if !strings.HasSuffix(info.Name(), TemplateFileExt) {
		t.Fatalf("TemplateFile does not end with %s", TemplateFileExt)
	}
}

func TestCheckFile(t *testing.T) {
	var err error
	// not existing file
	randomFileName := "hdljaskadssadsanlnk"
	if err = checkFile(randomFileName); err == nil {
		t.Fatalf("FAIL: checkFile(%s), exepected %v, actual %v", randomFileName, "error", err)
	} else {
		t.Logf("PASS: checkFile non existing file")
	}

	// directory as path
	tmpDirName, err := os.MkdirTemp("", "tmpdir")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer os.RemoveAll(tmpDirName)

	if err = checkFile(tmpDirName); err != ErrFoundDir {
		t.Fatalf("FAIL: checkFile(%s), exepected %v, actual %v", tmpDirName, ErrFoundDir, err)
	} else {
		t.Logf("PASS: checkFile on directory")
	}

	// valid tmpFile
	tmpFile, err := os.CreateTemp("", "tmpfile-")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	if err = checkFile(tmpFile.Name()); err != nil {
		t.Fatalf("FAIL: checkFile(%s), exepected %v, actual %v", tmpFile.Name(), nil, err)
	} else {
		t.Logf("PASS: checkFile on valid file")
	}
}

func TestParseFile(t *testing.T) {
	var err error
	tmpFile, err := os.CreateTemp("", "tmpfile-")
	if err != nil {
		t.Fatal(err.Error())
	}
	// defer tmpFile.Close()
	// defer os.Remove(tmpFile.Name())

	log.Printf("Temp File Name %s", tmpFile.Name())

	// valid input
	cmdList := []cmdItem{
		{
			Cmd:         "cmdddd",
			Description: "Descriptionnn",
		},
	}

	data, err := yaml.Marshal(&cmdList)
	if err != nil {
		t.Fatal(err.Error())
	}

	if _, err = tmpFile.Write(data); err != nil {
		t.Fatal(err.Error())
	}

	got, err := parseFile(tmpFile.Name())
	if err != nil || !reflect.DeepEqual(cmdList, got) {
		t.Fatalf("FAIL: parseFile(%s), expected(%v, %v), got(%v, %v)", tmpFile.Name(), cmdList, nil, got, err)
	}
	t.Logf("PASS: parseFile on valid file")

	// should not return the same cmdList
	got, _ = parseFile(tmpFile.Name())
	cmdList[0].Cmd = cmdList[0].Cmd + "123qwewqe"
	if reflect.DeepEqual(cmdList, got) {
		t.Fatalf("FAIL: parseFile(%s), expected(%v, %v), got(%v, %v)", tmpFile.Name(), cmdList, nil, got, err)
	}
	t.Logf("PASS: parseFile on valid file")
}

func TestGenBashScript(t *testing.T) {
	if err := GenBashScript(dummyTestFile); err != nil {
		log.Println(err.Error())
	}
}
func TestStartScenario(t *testing.T) {
	_ = StartScenario(dummyTestFile)
}
