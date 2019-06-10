package main

import (
	"fmt"
	"os"
	"time"

	"github.com/dheeraj-p/snapshot/targzhelper"
)

func createSnapshotDir(dirname string) error {
	if _, err := os.Stat(dirname); os.IsExist(err) {
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

func main() {
	var snapshotsDirName = ".snapshots"

	if err := createSnapshotDir(snapshotsDirName); err != nil {
		fmt.Println(err)
	}

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
}
