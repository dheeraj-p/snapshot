package main

import (
	"fmt"
	"os"

	"github.com/dheeraj-p/snapshot/targzhelper"
)

func main() {
	file, err := os.Create("snapshot.tar.gz")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var path = os.Args[1]

	err = targzhelper.MakeTar(path, file)

	if err != nil {
		fmt.Println(err.Error())
	}
}
