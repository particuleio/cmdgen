package cmdgen

import (
	"log"
	"os"
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

const (
	dummyTestFile = "../../files/test.yaml"
)

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
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	log.Printf("Temp File Name %s", tmpFile.Name())

	// valid input
	ts := templateStructure{
		Scenario: []cmdItem{
			{
				Cmd:         "cmdddd",
				Description: "Descriptionnn",
			},
		},
		Clean: []string{},
	}

	data, err := yaml.Marshal(&ts)
	if err != nil {
		t.Fatal(err.Error())
	}

	if _, err = tmpFile.Write(data); err != nil {
		t.Fatal(err.Error())
	}

	got, err := parseFile(tmpFile.Name())
	if err != nil || !reflect.DeepEqual(ts, got) {
		t.Fatalf("FAIL: parseFile(%s), expected(%v, %v), got(%v, %v)", tmpFile.Name(), ts, nil, got, err)
	}
	t.Logf("PASS: parseFile on valid file")

	// should not return the same cmdList
	got, _ = parseFile(tmpFile.Name())
	ts.Scenario[0].Cmd = ts.Scenario[0].Cmd + "123qwewqe"
	if reflect.DeepEqual(ts.Scenario, got) {
		t.Fatalf("FAIL: parseFile(%s), expected(%v, %v), got(%v, %v)", tmpFile.Name(), ts.Scenario, nil, got, err)
	}
	t.Logf("PASS: parseFile on valid file")
}

// func TestGenBashScript(t *testing.T) {
// 	if err := GenBashScript(dummyTestFile); err != nil {
// 		log.Println(err.Error())
// 	}
// }

func TestStartScenario(t *testing.T) {
	_ = StartScenario(dummyTestFile)
}

func TestCleanWorkspace(t *testing.T) {
	_ = CleanWorkspace(dummyTestFile)
}
