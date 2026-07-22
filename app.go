package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	goadb "github.com/electricbubble/gadb"
)

type App struct {
	currentPath string
	cache       DirCache
	ctx         context.Context
	client      goadb.Client
	device      goadb.Device
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) shutdown(_ context.Context) {
	if err := a.client.KillServer(); err != nil {
		log.Fatal(err)
	}
}

func (a *App) NewADBClient(adbPath string, port int) {
	if err := startADBServer(adbPath, port); err != nil {
		a.sendLogMsg(LogErr, err.Error())
		return
	}

	a.sendLogMsg(LogInfo, "adb server started successfully")
	client, err := goadb.NewClientWith("localhost", port)
	if err != nil {
		a.sendLogMsg(LogErr, err.Error())
		return
	}

	a.client = client
}

func (a *App) GetDeviceList() []string {
	devices, err := a.client.DeviceList()
	if err != nil {
		a.sendLogMsg(LogErr, err.Error())
		return nil
	}

	labels := make([]string, 0, len(devices))
	for _, device := range devices {
		info := device.DeviceInfo()
		serial := device.Serial()
		if info["device"] == "" {
			a.sendLogMsg(LogErr, serial+": device unauthorized")
			continue
		}

		// label := strings.Join([]string{serial, info["device"], info["model"]}, ":")
		labels = append(labels, serial)
	}

	return labels
}

func (a *App) SelectDevice(idx int) {
	devices, err := a.client.DeviceList()
	if err != nil {
		a.sendLogMsg(LogErr, err.Error())
		return
	}

	if l := len(devices); l <= 0 || l <= idx {
		a.sendLogMsg(LogErr, "Invalid device index")
		return
	}

	a.device = devices[idx]
	a.cache = *newDirCache(5)
}

func (a *App) DownloadADB() string {
	a.sendLogMsg(LogInfo, "downloading android platform tools")
	tmp, err := downloadADB()
	if err != nil {
		a.sendLogMsg(LogErr, err.Error())
		return ""
	}

	defer os.Remove(tmp)
	adbDir, err := getADBDir()
	if err != nil {
		a.sendLogMsg(LogErr, err.Error())
		return ""
	}

	a.sendLogMsg(LogInfo, "extracting android platform tools zip")
	if err := unzip(tmp, filepath.Dir(adbDir)); err != nil {
		a.sendLogMsg(LogErr, err.Error())
		return ""
	}

	adbPath, err := findADBExecutable(adbDir)
	if err != nil {
		a.sendLogMsg(LogErr, err.Error())
		return ""
	}

	return adbPath
}

func (a *App) KillServer(adbPath string, port int) {
	if err := killADBServer(adbPath, port); err != nil {
		log.Fatal(err)
	}
}
