package compress

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

type compress struct{}

// ICompress main iterface
type ICompress interface {
	Zip(string, string) (*os.File, error)
}

func (c *compress) Zip(filePath string, sourceFolder string) (*os.File, error) {
	zipFile, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}

	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)

	filepath.Walk(sourceFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		defer file.Close()

		writer, err := zipWriter.Create(info.Name())
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, file)
		if err != nil {
			return err
		}

		return nil
	})

	zipWriter.Close()

	return os.Open(filePath)
}

// New instance ...
func New() ICompress {
	return &compress{}
}
