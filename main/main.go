package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dheeraj-p/snapshot/targzhelper"
)

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

func takeSnapshot() {
	var snapshotsDirName = ".snapshots"
	destination := fmt.Sprintf("%s/snapshot_%s.tar.gz", snapshotsDirName, formattedTimeStamp())
	file, err := os.Create(destination)

	if err != nil {
		fmt.Println(err)
		return
	}

	var dirsToIgnore = []string{snapshotsDirName, ".git"}
	err = targzhelper.MakeTar("./", file, dirsToIgnore)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Snapshot is successfully taken")
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

func main() {
	setupSnapshotDirectory()

	if isNoOptionProvided() {
		showHelp()
		return
	}

	option := os.Args[1]

	if option == "take" {
		takeSnapshot()
		return
	}

	if option == "logs" {
		showLogs()
		return
	}

	showInvalidOption(option)
}
