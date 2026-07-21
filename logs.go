package main

import (
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type LogMsgKind int

const (
	LogErr LogMsgKind = iota
	LogWarn
	LogInfo
)

func (a *App) sendLogMsg(kind LogMsgKind, logs ...string) {
	var prefix rune
	switch kind {
	case LogErr:
		prefix = 'e'
	case LogWarn:
		prefix = 'w'
	case LogInfo:
		prefix = 'i'
	}

	logMsg := strings.Join([]string{string(prefix), ": ", strings.Join(logs, "")}, "")
	runtime.EventsEmit(a.ctx, "log-msg", logMsg)
}
