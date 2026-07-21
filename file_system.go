package main

import (
	"errors"
	"os"
	"path"
	"strings"

	goadb "github.com/electricbubble/gadb"
)

func (a *App) List(path string) []goadb.DeviceFileInfo {
	dirpath, err := cleanPath(path)
	if err != nil {
		a.sendLogMsg(LogErr, err.Error())
		return nil
	}

	if entries, ok := a.cache.get(dirpath); ok {
		a.currentPath = dirpath
		return entries
	}

	entries, err := a.device.List(dirpath)
	if err != nil {
		a.sendLogMsg(LogErr, err.Error())
		return nil
	}

	a.cache.set(dirpath, entries)
	a.currentPath = dirpath
	return entries
}

// TODO: add support for push/pull directories
// maybe compress and decompress directories

// NOTE: does not support pulling directories
func (a *App) Download(idx int, remote, local string) {
	files, ok := a.cache.get(a.currentPath)
	if !ok {
		a.sendLogMsg(LogErr, "current dir not found")
		return
	}

	if l := len(files); l <= 0 || l <= idx {
		a.sendLogMsg(LogErr, "invalid index")
		return
	}

	remoteDir, err := cleanPath(path.Dir(remote))
	if err != nil {
		a.sendLogMsg(LogErr, err.Error())
		return
	}

	if remote != path.Join(remoteDir, files[idx].Name) {
		// TODO: err message
		a.sendLogMsg(LogErr, "path error ", path.Join(remoteDir, files[idx].Name), remote)
		return
	}

	dest, err := os.OpenFile(local, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, files[idx].Mode)
	if err != nil {
		a.sendLogMsg(LogErr, err.Error())
		return
	}

	defer closeIO(dest)
	if err := a.device.Pull(remote, dest); err != nil {
		a.sendLogMsg(LogErr, err.Error())
		return
	}
}

// NOTE: local path should points to a file and remote path should points to a directory
func (a *App) Upload(local, remote string) {
	remoteDir, err := cleanPath(remote)
	if err != nil {
		a.sendLogMsg(LogErr, err.Error())
		return
	}

	file, err := os.Open(local)
	if err != nil {
		a.sendLogMsg(LogErr, err.Error())
		return
	}

	defer closeIO(file)
	stat, err := file.Stat()
	if err != nil {
		a.sendLogMsg(LogErr, err.Error())
		return
	}

	if err := a.device.Push(file, remoteDir, stat.ModTime(), stat.Mode()); err != nil {
		a.sendLogMsg(LogErr, err.Error())
		return
	}
	a.cache.invalidate(remoteDir)
}

func (a *App) Delete(path string) {
	remote, err := cleanPath(path)
	if err != nil {
		a.sendLogMsg(LogErr, err.Error())
		return
	}

	if _, err := a.device.RunShellCommand("rm -rf", remote); err != nil {
		a.sendLogMsg(LogErr, err.Error())
		return
	}
	a.cache.invalidateRec(remote)
}

func (a *App) Rename(old, new string) {
	oldPath, err := cleanPath(old)
	if err != nil {
		a.sendLogMsg(LogErr, err.Error())
		return
	}

	newPath, err := cleanPath(new)
	if err != nil {
		a.sendLogMsg(LogErr, err.Error())
		return
	}

	if oldp := path.Dir(oldPath); oldp != path.Dir(newPath) || !strings.HasPrefix(oldp, a.currentPath) {
		a.sendLogMsg(LogErr, "cannot move")
		return
	}

	if _, err := a.device.RunShellCommand("mv", oldPath, newPath); err != nil {
		a.sendLogMsg(LogErr, err.Error())
		return
	}
	a.cache.invalidateRec(a.currentPath)
}

func (a *App) MakeDir(dirname string) {
	dirPath, err := cleanPath(path.Join(a.currentPath, dirname))
	if err != nil {
		a.sendLogMsg(LogErr, err.Error())
		return
	}

	if !strings.HasPrefix(dirPath, a.currentPath) {
		a.sendLogMsg(LogErr, "not in current dir")
		return
	}

	if _, err := a.device.RunShellCommand("mkdir", dirPath); err != nil {
		a.sendLogMsg(LogErr, err.Error())
		return
	}
	a.cache.invalidate(a.currentPath)
}

func cleanPath(fpath string) (string, error) {
	const STORAGE_DIR = "/storage/"

	if strings.ContainsAny(";&|", fpath) {
		return "", errors.New("invalid characters")
	}

	cleanPath := path.Clean(fpath)
	if cleanPath != STORAGE_DIR[:len(STORAGE_DIR)-1] && !strings.HasPrefix(cleanPath, STORAGE_DIR) {
		return "", errors.New("path escapes " + STORAGE_DIR)
	}

	return cleanPath, nil
}
