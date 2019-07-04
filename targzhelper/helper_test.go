package targzhelper

import (
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
	testFilePath := fmt.Sprintf("%s/file%s.tar.gz", testDirPath, time.Now())
	file, _ := os.Create(testFilePath)
	invalidDirPath := filepath.Join(testDirPath, "invalidPath")

	err := MakeTar(invalidDirPath, file, []string{})
	expected := fmt.Errorf("Invalid Path stat %s: no such file or directory", invalidDirPath)

	if !reflect.DeepEqual(err, expected) {
		t.Errorf("\nEXPECTED === %s\nGOT === %s\n", expected, err)
	}
}

func testMakeTarWithInvalidPermission(t *testing.T, testDirPath string) {
	testFilePath := fmt.Sprintf("%s/file%s.tar.gz", testDirPath, time.Now())
	file, _ := os.Create(testFilePath)
	testPath := filepath.Join(testDirPath, "dir")
	os.Mkdir(testPath, 0000)

	err := MakeTar(testPath, file, []string{})
	expected := fmt.Errorf("during walk open %s: permission denied", testPath)

	if !reflect.DeepEqual(err, expected) {
		t.Errorf("\nEXPECTED === %s\nGOT === %s\n", expected, err)
	}
}

func testMakeTarWithFile(t *testing.T, testDirPath string) {
	testFilePath := fmt.Sprintf("%s/file%s.tar.gz", testDirPath, time.Now())
	file, _ := os.Create(testFilePath)
	testPath := filepath.Join(testDirPath, "file")
	os.Create(testPath)

	err := MakeTar(testPath, file, []string{})
	expected := fmt.Errorf("%s is not a directory", testPath)

	if !reflect.DeepEqual(err, expected) {
		t.Errorf("\nEXPECTED === %s\nGOT === %s\n", expected, err)
	}
}

func testMakeTarWithValidArgs(t *testing.T, testDirPath string) {
	testFilePath := fmt.Sprintf("%s/file%s.tar.gz", testDirPath, time.Now())
	file, _ := os.Create(testFilePath)
	testPath := filepath.Join(testDirPath, "validDir")
	os.Mkdir(testPath, 0777)

	err := MakeTar(testPath, file, []string{})
	if err != nil {
		t.Errorf("Unexpected error occured %s", err.Error())
	}
}

func TestMakeTar(t *testing.T) {
	snapshotDir, err := setUpForMakeTar()

	if err != nil {
		t.Errorf("Could not setup tests %s", err.Error())
		return
	}

	testMakeTarWithInvalidPath(t, snapshotDir)
	testMakeTarWithFile(t, snapshotDir)
	testMakeTarWithValidArgs(t, snapshotDir)
	testMakeTarWithInvalidPermission(t, snapshotDir)

	if err := tearDownForMakeTar(snapshotDir); err != nil {
		t.Errorf("Could not tearDown tests")
	}
}

func setUpForUntar() (string, string) {
	testDirPath, _ := setUpForMakeTar()
	testFilePath := fmt.Sprintf("%s/file%s.tar.gz", testDirPath, time.Now())
	file, _ := os.Create(testFilePath)
	testPath := filepath.Join(testDirPath, "validDir")
	os.Mkdir(testPath, 0777)

	MakeTar(testPath, file, []string{})
	return testFilePath, testDirPath
}

func testUntarWithValidArgs(t *testing.T, testDirPath, tarFilePath string) {
	file, _ := os.OpenFile(tarFilePath, os.O_RDONLY, 0555)

	err := Untar(file, testDirPath)

	if err != nil {
		t.Errorf("Unexpected error happened %s", err)
	}
}

func testUntarWithDirectory(t *testing.T, testDirPath string) {
	dir := os.TempDir()
	path, _ := os.OpenFile(dir, os.O_RDONLY, 0555)
	defer func() {
		os.RemoveAll(dir)
	}()

	err := Untar(path, testDirPath)
	expected := fmt.Sprintf("read %s: is a directory", path.Name())

	if !reflect.DeepEqual(err.Error(), expected) {
		t.Errorf("\nEXPECTED === %s\nGOT === %s\n", expected, err)
	}
}

func TestUntar(t *testing.T) {
	tarFilePath, testDirPath := setUpForUntar()

	testUntarWithValidArgs(t, testDirPath, tarFilePath)
	testUntarWithDirectory(t, testDirPath)
}
