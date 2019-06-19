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

func extractFile(header tar.Header, reader io.Reader, path string) (rerr error) {
	fileName := filepath.Join(path, header.Name)
	file, ferr := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, header.FileInfo().Mode())

	defer func() {
		err := file.Close()
		fmt.Errorf("error while closing the file %s", err)
	}()

	if ferr != nil {
		return ferr
	}

	_, err := io.Copy(file, reader)

	if err != nil {
		return err
	}

	return nil
}

func extractDir(header tar.Header, reader io.Reader, path string) error {
	dirName := filepath.Join(path, header.Name)
	err := os.MkdirAll(dirName, header.FileInfo().Mode())

	if err != nil {
		return err
	}
	return nil
}

func Untar(reader io.Reader, path string) error {
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}

	defer func() {
		err := gzipReader.Close()
		fmt.Errorf("error closing the reader %s", err)
	}()
	tarReader := tar.NewReader(gzipReader)

	for header, err := tarReader.Next(); err != io.EOF; header, err = tarReader.Next() {
		if err != nil {
			return err
		}
		extract := extractFile
		if header.FileInfo().IsDir() {
			extract = extractDir
		}

		if err := extract(*header, tarReader, path); err != nil {
			return err
		}
	}

	return nil
}

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

func MakeTar(path string, writer io.Writer, pathToIgnore []string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("Invalid Path %s", err.Error())
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
			return fmt.Errorf("during walk %s", err)
		}

		header, err := tar.FileInfoHeader(fileInfo, fileInfo.Name())
		if err != nil {
			return fmt.Errorf("invalid file header %s", err)
		}

		header.Name = strings.TrimPrefix(strings.Replace(fileName, path, "", -1), string(filepath.Separator))

		if contains(pathToIgnore, header.Name) {
			return nil
		}

		if hasParentIn(pathToIgnore, header.Name) {
			return nil
		}

		if header.Name == "" {
			return nil
		}

		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		file, err := os.Open(fileName)
		defer func() {
			err := file.Close()
			fmt.Errorf("error while closing the file %s", err)
		}()

		if err != nil {
			return err
		}

		if _, err := io.Copy(tarWriter, file); err != nil {
			return err
		}

		return nil
	})
}
