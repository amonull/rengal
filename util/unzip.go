package util

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/amonull/rengal/filesystem"
)

func Unzip(zipStream io.ReaderAt, size int64, dest string) error {
	r, err := zip.NewReader(zipStream, size)
	if err != nil {
		return err
	}

	err = filesystem.Api().MkdirAll(dest, os.ModePerm)
	if err != nil {
		return err
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f, dest)

		if err != nil {
			return err
		}
	}

	return nil
}

// Closure to address file descriptors issue with all the deferred .Close() methods
func extractAndWriteFile(zipFile *zip.File, dest string) error {
	rc, err := zipFile.Open()
	if err != nil {
		return err
	}

	defer Ignore(rc.Close)

	path := filepath.Join(dest, zipFile.Name)

	// Check for ZipSlip (Directory traversal)
	if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
		return fmt.Errorf("illegal file path: %s", path)
	}

	if zipFile.FileInfo().IsDir() {
		err = filesystem.Api().MkdirAll(path, zipFile.Mode())
		if err != nil {
			return err
		}
	}

	err = filesystem.Api().MkdirAll(filepath.Dir(path), zipFile.Mode())
	if err != nil {
		return err
	}

	f, err := filesystem.Api().OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zipFile.Mode())
	if err != nil {
		return err
	}

	defer Ignore(f.Close)

	_, err = io.Copy(f, rc)
	if err != nil {
		return err
	}

	return nil
}
