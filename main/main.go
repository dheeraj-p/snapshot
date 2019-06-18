package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dheeraj-p/snapshot/targzhelper"
)

var snapshotMessages map[string]string

func createDirectoryIfNotExists(dirname string) error {
	if _, err := os.Stat(dirname); err == nil {
		return nil
	}

	if err := os.Mkdir(dirname, 0777); err != nil {
		return err
	}

	return nil
}

func formattedTimeStamp() string {
	return time.Now().Format("2006_01_02_15_04_05")
}

func setupSnapshotDirectory() {
	if err := createDirectoryIfNotExists(".snapshots"); err != nil {
		fmt.Println("Came Here")
		fmt.Println(err)
	}
}

func takeSnapshot() (string, error) {
	var snapshotsDirName = ".snapshots"

	if len(os.Args) < 3 {
		return "", fmt.Errorf("not enough arguments")
	}

	destination := fmt.Sprintf("%s/snapshot_%s.tar.gz", snapshotsDirName, formattedTimeStamp())
	snapshotMessages[destination] = os.Args[2]

	file, err := os.Create(destination)

	if err != nil {
		return "", err
	}

	var dirsToIgnore = []string{snapshotsDirName, ".git"}
	err = targzhelper.MakeTar("./", file, dirsToIgnore)
	if err != nil {
		return "", err
	}

	return "Snapshot is successfully taken", nil
}

func showLogs() {
	filepath.Walk("./.snapshots", func(fileName string, fileInfo os.FileInfo, err error) error {
		if !fileInfo.IsDir() {
			fmt.Println(strings.TrimPrefix(fileName, ".snapshots/"))
		}
		return nil
	})
}

func showInvalidOption(option string) {
	fmt.Printf("Invalid Option - %s\n", option)
}

func showHelp() {
	var takeSnapshotHelp = "option: take - To take snapshot of current state"
	var showLogsHelp = "option: logs - Show info about all the snapshots taken"
	fmt.Printf("%s\n%s\n", takeSnapshotHelp, showLogsHelp)
}

func isNoOptionProvided() bool {
	return len(os.Args) < 2
}

func writeToFile() {
	buffer, _ := json.Marshal(snapshotMessages)
	f, _ := os.OpenFile(".snapshots/data.json", os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	f.Write(buffer)
}

func main() {
	setupSnapshotDirectory()

	snapshotMessages = make(map[string]string)
	buffer, _ := ioutil.ReadFile(".snapshots/data.json")
	json.Unmarshal(buffer, &snapshotMessages)

	if isNoOptionProvided() {
		showHelp()
		return
	}

	option := os.Args[1]

	if option == "take" {
		str, err := takeSnapshot()
		writeToFile()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(str)
		return
	}

	if option == "logs" {
		showLogs()
		return
	}

	showInvalidOption(option)
}
