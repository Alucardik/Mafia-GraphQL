package utils

import "log"

func NotifyOnError(msg string, err error) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}

func PanicOnError(msg string, err error) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func FailOnError(msg string, err error) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
