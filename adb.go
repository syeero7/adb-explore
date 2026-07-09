package main

import (
	"archive/zip"
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const ADB_DECOMPRESSED_MAX_SIZE = 1024 * 1024 * 30
const OFFICIAL_ADB_PAGE_URL = "https://developer.android.com/tools/releases/platform-tools"
const ADB_DOWNLOAD_URL_PREFIX = "https://dl.google.com/android/repository/platform-tools-latest-"

var ErrSizeLimitExceeded = errors.New("size limit exceeded")

var adbExePath string

// https://dl.google.com/android/repository/platform-tools-latest-windows.zip
// https://dl.google.com/android/repository/platform-tools-latest-darwin.zip
// https://dl.google.com/android/repository/platform-tools-latest-linux.zip

func downloadADB() (string, error) {
	parts := strings.Split(runtime.GOOS, "/")
	res, err := http.Get(ADB_DOWNLOAD_URL_PREFIX + parts[0] + ".zip")
	if err != nil {
		return "", err
	}
	defer closeIO(res.Body)

	file, err := os.CreateTemp("", "platform-tools.zip")
	if err != nil {
		return "", err
	}
	defer closeIO(file)

	_, err = io.Copy(file, res.Body)
	return file.Name(), err
}

func startADBServer() error {
	var err0 error
	adbExePath, err0 = getADBPath()
	if err0 != nil {
		tmp, err := downloadADB()
		if err == nil {
			return err
		}

		defer os.Remove(tmp)
		adbDir, err := getADBDir()
		if err != nil {
			return err
		}

		if err := unzip(tmp, adbDir); err != nil {
			return err
		}

		adbExePath, err = findADBExecutable(adbDir)
		if err != nil {
			return err
		}
	}

	return exec.Command(adbExePath, "start-server").Run()
}

func killADBServer() error {
	return exec.Command(adbExePath, "kill-server").Run()
}

func getADBDir() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	return filepath.Join(filepath.Dir(exePath), "platform-tools"), nil
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer closeIO(r)

	var bytesRead int64
	for _, file := range r.File {
		fname, err := filepath.Localize(file.Name)
		if err != nil {
			return err
		}

		path := filepath.Join(dest, fname)

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(path, file.Mode()); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return err
		}

		if err := extractFile(file, path, &bytesRead); err != nil {
			return err
		}
	}

	return nil
}

func extractFile(file *zip.File, dest string, bytesRead *int64) error {
	remaining := ADB_DECOMPRESSED_MAX_SIZE - *bytesRead
	if remaining <= 0 {
		return ErrSizeLimitExceeded
	}

	srcF, err := file.Open()
	if err != nil {
		return err
	}
	defer closeIO(srcF)

	destF, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return err
	}
	defer closeIO(destF)

	n, err := io.Copy(destF, io.LimitReader(srcF, remaining))
	if err != nil {
		return err
	}

	*bytesRead += n
	if *bytesRead >= ADB_DECOMPRESSED_MAX_SIZE {
		return ErrSizeLimitExceeded
	}

	return nil
}

func getADBPath() (string, error) {
	path, err := exec.LookPath("adb")
	if err == nil {
		return path, nil
	}

	adbDir, err := getADBDir()
	if err != nil {
		return "", err
	}

	return findADBExecutable(adbDir)
}

func findADBExecutable(adbDir string) (string, error) {
	entries, err := os.ReadDir(adbDir)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "adb") && !entry.IsDir() {
			return filepath.Join(adbDir, entry.Name()), nil
		}
	}

	return "", os.ErrNotExist
}
