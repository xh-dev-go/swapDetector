package main

import (
	"encoding/json"
	"flag"
	"os"
	"time"
)

const VERSION = "1.0.0"

type UpdateType string

const (
	Swap   UpdateType = "SWAP"
)

type UpdateMsg struct {
	FilePath string
	Type     UpdateType
}

func (updateMsg *UpdateMsg) ToString() string {
	if bytes, err := json.Marshal(updateMsg); err != nil {
		panic(err)
	} else {
		return string(bytes)
	}
}

func detect(
	file string, update chan UpdateMsg, delay int64) {
	var lastSize int64

	for {
		if fStat, err := os.Stat(file); err != nil {
			panic(err)
		} else {
			size := fStat.Size()

			//mod := fStat.ModTime()

			if size < lastSize {
				update <- UpdateMsg{
					FilePath: file,
					Type:     Swap,
				}
				break
			}

			if lastSize < size {
				lastSize = size
			}

			if delay > 0 {
				duration := time.Duration(delay) * time.Millisecond
				time.Sleep(duration)
			}
		}
	}
}

func main() {
	const cmdVersion = "version"
	const cmdPath = "file"
	const cmdDelay = "delay"
	var file string
	var delay int64
	var checkVersion bool
	flag.BoolVar(&checkVersion, cmdVersion, false, "check the version fo application")
	flag.StringVar(&file, cmdPath, "", "the path to detect")
	flag.Int64Var(&delay, cmdDelay, -1, "the time delay for check file stat")
	flag.Parse()

	if checkVersion {
		println(VERSION)
		os.Exit(0)
	}

	if file == "" {
		print("Please input file path")
		flag.Usage()
		os.Exit(1)
	}

	if delay == -1 {
		delay = 1000
	}

	fileUpdate := make(chan UpdateMsg)

	go detect(file, fileUpdate, delay)

	select {
	case msg := <-fileUpdate:
		println(msg.FilePath)
	}
}
