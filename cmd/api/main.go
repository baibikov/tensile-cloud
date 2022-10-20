package main

import "github.com/sirupsen/logrus"

func main() {
	if err := runApp(); err != nil {
		logrus.Fatal(err)
	}
}
