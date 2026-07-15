package main

import (
	"context"
	"log"

	goadb "github.com/electricbubble/gadb"
)

type App struct {
	ctx    context.Context
	client goadb.Client
	FileSystem
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) NewADBClient() []string {
	if err := startADBServer(); err != nil {
		log.Fatal(err)
	}

	client, err := goadb.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	a.client = client
	devices, err := client.DeviceList()
	if err != nil {
		log.Fatal(err)
	}

	labels := make([]string, 0, len(devices))
	for _, device := range devices {
		info := device.DeviceInfo()
		label := info["device"] + ":" + info["model"]
		labels = append(labels, label)
	}

	return labels
}

func (a *App) SelectDevice(idx int) {
	devices, err := a.client.DeviceList()
	if err != nil {
		log.Fatal(err)
	}

	if l := len(devices); l <= 0 || l <= idx {
		log.Fatal("Invalid device index")
	}

	a.FileSystem.init(&devices[idx])
}
