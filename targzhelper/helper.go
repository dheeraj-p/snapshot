package targzhelper

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func contains(collection []string, text string) bool {
	for index := range collection {
		if collection[index] == text {
			return true
		}
	}
	return false
}

func hasParentIn(parentCandidates []string, childPath string) bool {
	for index := range parentCandidates {
		if isParentOf(parentCandidates[index], childPath) {
			return true
		}
	}
	return false
}

func isParentOf(parentDir string, pathToCheck string) bool {
	return strings.HasPrefix(pathToCheck, parentDir+"/")
}

func MakeTar(path string, writer io.Writer, dirsToIgnore []string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !fileInfo.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}

	gzWriter := gzip.NewWriter(writer)
	defer gzWriter.Close()

	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	return filepath.Walk(path, func(fileName string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(fileInfo, fileInfo.Name())
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(strings.Replace(fileName, path, "", -1), string(filepath.Separator))

		if contains(dirsToIgnore, header.Name) {
			return nil
		}

		if hasParentIn(dirsToIgnore, header.Name) {
			return nil
		}

		if header.Name == "" {
			return nil
		}

		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		if !fileInfo.Mode().IsRegular() {
			return nil
		}

		file, err := os.Open(fileName)

		if err != nil {
			return err
		}

		if _, err := io.Copy(tarWriter, file); err != nil {
			return err
		}

		file.Close()
		return nil
	})
}
