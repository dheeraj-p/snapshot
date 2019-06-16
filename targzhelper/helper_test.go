package targzhelper

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func setUpForMakeTar() (string, error) {
	tempDir := os.TempDir()
	timeStamp := time.Now().Format("2006_01_02_15_04_05")
	snapshotTestDir := filepath.Join(tempDir, timeStamp)
	if err := os.Mkdir(snapshotTestDir, 0777); err != nil {
		return "", err
	}
	return snapshotTestDir, nil
}

func tearDownForMakeTar(dirPath string) error {
	err := os.RemoveAll(dirPath)
	return err
}

func testMakeTarWithInvalidPath(t *testing.T, testDirPath string) {
	buffer := bytes.NewBuffer([]byte{})
	invalidDirPath := filepath.Join(testDirPath, "foo")
	err := MakeTar(invalidDirPath, buffer, []string{})
	expected := fmt.Errorf("Invalid Path stat %s: no such file or directory", invalidDirPath)
	if !reflect.DeepEqual(err, expected) {
		t.Errorf("\nEXPECTED === %s\nGOT === %s\n", expected, err)
	}
}

func testMakeTarWithFile(t *testing.T, testDirPath string) {
	buffer := bytes.NewBuffer([]byte{})
	testFilePath := filepath.Join(testDirPath, "bar")
	os.Create(testFilePath)
	err := MakeTar(testFilePath, buffer, []string{})
	expected := fmt.Errorf("%s is not a directory", testFilePath)
	if !reflect.DeepEqual(err, expected) {
		t.Errorf("\nEXPECTED === %s\nGOT === %s\n", expected, err)
	}
}

func TestMakeTar(t *testing.T) {
	os.Mkdir("tmp", 0777)
	snapshotDir, err := setUpForMakeTar()

	if err != nil {
		t.Errorf("Could not setup tests")
		return
	}

	testMakeTarWithInvalidPath(t, snapshotDir)
	testMakeTarWithFile(t, snapshotDir)

	if err := tearDownForMakeTar(snapshotDir); err != nil {
		t.Errorf("Could not tearDown tests")
	}
}
