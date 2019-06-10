package main

import (
	"fmt"
	"os"
	"time"

	"github.com/dheeraj-p/snapshot/targzhelper"
)

func createSnapshotDir(dirname string) error {
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		if err = os.Mkdir(dirname, 0777); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	timeStamp := time.Now().Format("2006_01_02_15_04_05")

	if err := createSnapshotDir(".snapshots"); err != nil {
		fmt.Println(err)
	}

	destination := fmt.Sprintf(".snapshots/snapshot_%s.tar.gz", timeStamp)
	file, err := os.Create(destination)

	if err != nil {
		fmt.Println(err)
		return
	}

	err = targzhelper.MakeTar("./", file, []string{".snapshots", ".git"})
	if err != nil {
		fmt.Println(err)
	}
}
