package main

import (
	"cmp"
	"errors"
	"os"
	"path"
	"slices"
	"strings"
	"time"
)

type Entry struct {
	IsDir        bool        `json:"isDir"`
	Name         string      `json:"name"`
	Path         string      `json:"path"`
	Size         string      `json:"size"`
	Mode         os.FileMode `json:"mode"`
	LastModified time.Time   `json:"lastModified"`
}

type DirEntries struct {
	Parent  string  `json:"parent"`
	Path    string  `json:"path"`
	Entries []Entry `json:"entries"`
}

func (a *App) List(path, query, sortBy string) DirEntries {
	dirpath, err := cleanPath(path)
	if err != nil {
		a.sendLogMsg(LogErr, err.Error())
		return DirEntries{}
	}

	if entries, ok := a.cache.get(dirpath); ok {
		a.currentPath = dirpath
		return a.sortFilterDir(&entries, query, sortBy)
	}

	entries, err := a.getEntries(dirpath)
	if err != nil {
		a.sendLogMsg(LogErr, err.Error())
		return DirEntries{}
	}

	a.currentPath = dirpath
	return a.sortFilterDir(&entries, query, sortBy)
}

// TODO: add support for push/pull directories
// maybe compress and decompress directories

// NOTE: does not support pulling directories
func (a *App) Download(idx int, remote, local string) {
	entries, ok := a.cache.get(a.currentPath)
	if !ok {
		a.sendLogMsg(LogErr, "current dir not found")
		return
	}

	if l := len(entries.Entries); l <= 0 || l <= idx {
		a.sendLogMsg(LogErr, "invalid index")
		return
	}

	remoteDir, err := cleanPath(path.Dir(remote))
	if err != nil {
		a.sendLogMsg(LogErr, err.Error())
		return
	}

	if remote != path.Join(remoteDir, entries.Entries[idx].Name) {
		// TODO: err message
		a.sendLogMsg(LogErr, "path error ", path.Join(remoteDir, entries.Entries[idx].Name), remote)
		return
	}

	dest, err := os.OpenFile(local, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, entries.Entries[idx].Mode)
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

func (a *App) getEntries(dirpath string) (DirEntries, error) {
	items, err := a.device.List(dirpath)
	if err != nil {
		return DirEntries{}, err
	}

	entries := make([]Entry, 0, max(0, len(items)-2))
	for _, item := range items {
		if item.Name == "." || item.Name == ".." {
			continue
		}

		entry := Entry{
			IsDir:        item.IsDir(),
			Name:         item.Name,
			Mode:         item.Mode,
			LastModified: item.LastModified,
			Path:         path.Join(dirpath, item.Name),
			// TODO: Size: item
		}

		if entry.IsDir && !strings.HasSuffix(entry.Path, "/") {
			entry.Path += "/"
		}

		entries = append(entries, entry)
	}

	return DirEntries{Path: dirpath, Parent: path.Dir(dirpath), Entries: entries}, nil
}

func (a *App) sortFilterDir(dir *DirEntries, query, sortBy string) DirEntries {
	sorted := sortEntries(dir, sortBy)
	entries, filtered := filterEntries(sorted, strings.TrimSpace(query))
	a.cache.set(a.currentPath, *entries)
	entries.Entries = filtered
	return *entries
}

func filterEntries(dir *DirEntries, query string) (*DirEntries, []Entry) {
	if query == "" {
		return dir, dir.Entries
	}

	entries := make([]Entry, 0, len(dir.Entries))
	for _, entry := range dir.Entries {
		if strings.Contains(entry.Name, query) {
			entries = append(entries, entry)
		}
	}

	return dir, entries
}

func sortEntries(dir *DirEntries, sortBy string) *DirEntries {
	parts := strings.Split(sortBy, ":")
	slices.SortFunc(dir.Entries, func(a, b Entry) int {
		switch parts[0] {
		case "name":
			if parts[1] == "asc" {
				return cmp.Compare(a.Name, b.Name)
			}

			return cmp.Compare(b.Name, a.Name)
		case "size":
			if parts[1] == "asc" {
				return cmp.Compare(a.Size, b.Size)
			}

			return cmp.Compare(b.Size, a.Size)
		default:
			if parts[1] == "asc" {
				return a.LastModified.Compare(b.LastModified)
			}

			return b.LastModified.Compare(a.LastModified)
		}
	})

	return dir
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
