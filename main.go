package main

import (
	"os"

	"github.com/holi0317/bb2gh/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	app := cmd.New()
	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
