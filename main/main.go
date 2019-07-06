package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dheeraj-p/snapshot/targzhelper"
)

type snapshot struct {
	Message   string
	Timestamp int64
	FileName  string
}

var snapshots map[string]snapshot

func logError(err error) {
	log.Fatal(err)
}

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

func setupSnapshotDirectory(path string) error {
	if err := createDirectoryIfNotExists(path + "/.snapshots"); err != nil {
		return err
	}
	return nil
}

func getpathToIgnore(basePath string) []string {
	buffer, err := ioutil.ReadFile(basePath + "/.signore")
	if err != nil {
		return []string{}
	}
	fileContent := strings.Trim(string(buffer), "\n")
	pathToIgnore := strings.Split(fileContent, "\n")
	return pathToIgnore
}

func takeSnapshot(snapshotsDirBase string) (string, error) {
	var snapshotsDirName = snapshotsDirBase + "/.snapshots"

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

	pathToIgnore := getpathToIgnore(snapshotsDirBase)

	err = targzhelper.MakeTar(snapshotsDirBase, file, pathToIgnore)
	if err != nil {
		return "", err
	}

	return "Snapshot is successfully taken", nil
}

func formatLog(sn snapshot, sha string) string {
	timestamp := formatTimeStamp(sn.Timestamp, "Mon Jan _2 15:04:05 2006")
	return fmt.Sprintf("Snapshot Id:\t%s\nDate:\t%s\n\n\t%s\n", sha, timestamp, sn.Message)
}

func validate(path string) string {
	hasSnapshotDir := false
	files, _ := ioutil.ReadDir(path)

	for _, f := range files {
		if f.IsDir() && f.Name() == ".snapshots" {
			hasSnapshotDir = true
		}
	}

	if !hasSnapshotDir {
		paths := strings.Split(path, "/")
		basePath := strings.Join(paths[:len(paths)-1], "/")
		if basePath == "/Users" {
			log.Fatal("This is not a snapshot directory. Use snapshot init to initialize snapshot")
		}
		return validate(basePath)
	}
	return path
}

func showLogs() {
	for sha := range snapshots {
		_snapshot := snapshots[sha]
		log := formatLog(_snapshot, sha)
		fmt.Println(log)
	}
}

func checkout(snapshotBaseDir string) error {
	if len(os.Args) < 4 {
		return fmt.Errorf("not enough arguments")
	}

	sha := os.Args[2]
	directory := os.Args[3]
	fileName := snapshots[sha].FileName

	file, err := os.OpenFile(fileName, os.O_RDONLY, 0777)
	if err != nil {
		return err
	}

	dirname := snapshotBaseDir + "/" + directory + "/snapshot_" + sha

	err = os.MkdirAll(dirname, 0777)
	if err != nil {
		return err
	}

	err = targzhelper.Untar(file, dirname)
	if err != nil {
		return err
	}

	fmt.Printf("Checked out version is available in ---> %s\n", dirname)
	return nil
}

func showInvalidOption(option string) {
	fmt.Printf("Invalid Option - %s\n", option)
}

func showHelp() {
	takeSnapshotHelp := "option: take [message] - To take snapshot of current state"
	showLogsHelp := "option: logs - Show info about all the snapshots taken"
	checkoutHelp := "option: checkout [snapshot id] [directory to checkout] - Checkout a specific snapshot by providing snapshot id"
	initHelp := "option: init- Initialize directory as snapshot directory"
	note := "NOTE: Snapshot id can be found under logs"
	fmt.Printf("%s\n\n%s\n\n%s\n\n%s\n\n%s\n", takeSnapshotHelp, showLogsHelp, checkoutHelp, initHelp, note)
}

func isNoOptionProvided() bool {
	return len(os.Args) < 2
}

func writeToFile(snapshotDir string) {
	buffer, _ := json.Marshal(snapshots)
	f, _ := os.OpenFile(snapshotDir+"/.snapshots/data.json", os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	_, err := f.Write(buffer)
	if err!=nil{
		log.Fatalf("Error occured while writing to data.json\n %s", err.Error())
	}
}

func createDataIfNotExists(fileName string) {
	if _, err := os.Stat(fileName); err != nil {
		file, _ := os.Create(fileName)
		_, err := file.WriteString("{}")
		if err!=nil{
			log.Fatalf("Error occured initializing data.json\n %s", err.Error())
		}
	}
}

func main() {
	if isNoOptionProvided() {
		showHelp()
		return
	}

	option := os.Args[1]
	currentPath, _ := os.Getwd()

	if option == "init" {
		err := setupSnapshotDirectory(currentPath)
		if err != nil {
			logError(err)
		}
		createDataIfNotExists(currentPath + "/.snapshots/data.json")

		signoreFilePath := currentPath + "/.signore"
		file, _ := os.Create(signoreFilePath)
		_, err = file.WriteString(".snapshots\ncheckedout_versions")
		if err!=nil{
			log.Fatalf("Error occured while writing to .signore\n %s", err.Error())
		}
		
		fmt.Print(currentPath + "is initialized as snapshot directory")
		return
	}

	snapshotDir := validate(currentPath)

	dataFilePath := snapshotDir + "/.snapshots/data.json"
	snapshots = make(map[string]snapshot)

	buffer, err := ioutil.ReadFile(dataFilePath)
	if err != nil {
		logError(err)
	}

	err = json.Unmarshal(buffer, &snapshots)
	if err != nil {
		logError(err)
	}

	if option == "take" {
		str, err := takeSnapshot(snapshotDir)
		writeToFile(snapshotDir)
		if err != nil {
			logError(err)
		}
		fmt.Println(str)
		return
	}

	if option == "logs" {
		showLogs()
		return
	}

	if option == "checkout" {
		err := checkout(snapshotDir)
		if err != nil {
			logError(err)
		}
		return
	}

	showInvalidOption(option)
}
