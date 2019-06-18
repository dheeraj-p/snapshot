package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/dheeraj-p/snapshot/targzhelper"
)

type snapshot struct {
	Message   string
	Timestamp int64
	FileName  string
}

var snapshots map[string]snapshot

func createDirectoryIfNotExists(dirname string) error {
	if _, err := os.Stat(dirname); err == nil {
		return nil
	}

	if err := os.Mkdir(dirname, 0777); err != nil {
		return err
	}

	return nil
}

func formatTimeStamp(timestamp int64, format string) string {
	_time := time.Unix(timestamp, 0)
	return _time.Format(format)
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

	timestamp := time.Now().Unix()
	sha := fmt.Sprintf("%x", timestamp)
	message := os.Args[2]

	formattedTimeStamp := formatTimeStamp(timestamp, "2006_01_02_15_04_05")
	destination := fmt.Sprintf("%s/snapshot_%s.tar.gz", snapshotsDirName, formattedTimeStamp)

	snapshots[sha] = snapshot{message, timestamp, destination}

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

func formatLog(sn snapshot, sha string) string {
	timestamp := formatTimeStamp(sn.Timestamp, "Mon Jan _2 15:04:05 2006")
	return fmt.Sprintf("Snapshot Id:\t%s\nDate:\t%s\n\n\t%s\n", sha, timestamp, sn.Message)
}

func showLogs() {
	for sha := range snapshots {
		_snapshot := snapshots[sha]
		log := formatLog(_snapshot, sha)
		fmt.Println(log)
	}
}

func checkout() {
	if len(os.Args) < 3 {
		fmt.Errorf("not enough arguments")
	}

	sha := os.Args[2]
	fileName := snapshots[sha].FileName

	file, _ := os.OpenFile(fileName, os.O_RDONLY, 0777)
	dirname := "checkedout_versions/snapshot_" + sha

	os.MkdirAll(dirname, 0777)
	targzhelper.Untar(file, dirname)

	fmt.Printf("Checked out version is available in ---> %s\n", dirname)
}

func showInvalidOption(option string) {
	fmt.Printf("Invalid Option - %s\n", option)
}

func showHelp() {
	takeSnapshotHelp := "option: take [message] - To take snapshot of current state"
	showLogsHelp := "option: logs - Show info about all the snapshots taken"
	checkoutHelp := "option: checkout [snapshot id]- Checkout a specific snapshot by providing snapshot id"
	note := "NOTE: Snapshot id can be found under logs"
	fmt.Printf("%s\n%s\n%s\n\n%s\n", takeSnapshotHelp, showLogsHelp, checkoutHelp, note)
}

func isNoOptionProvided() bool {
	return len(os.Args) < 2
}

func writeToFile() {
	buffer, _ := json.Marshal(snapshots)
	f, _ := os.OpenFile(".snapshots/data.json", os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	f.Write(buffer)
}

func main() {
	setupSnapshotDirectory()

	snapshots = make(map[string]snapshot)
	buffer, _ := ioutil.ReadFile(".snapshots/data.json")
	json.Unmarshal(buffer, &snapshots)

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

	if option == "checkout" {
		checkout()
		return
	}

	showInvalidOption(option)
}
