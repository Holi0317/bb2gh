package main

import (
	"io"
	"os"

	"github.com/holi0317/bb2gh/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	file, err := os.OpenFile("bb2gh.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic("Failed to open log file bb2gh.log")
	}

	defer file.Close()

	logrus.SetOutput(io.MultiWriter(os.Stdout, file))

	app := cmd.New()
	err = app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
