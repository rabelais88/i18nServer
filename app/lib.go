package app

import (
	"log"
	"os"
)

func checkError(e error, msg string) {
	if e != nil {
		if msg != "" {
			log.Print(msg)
		}
		log.Print(e)
	}
}

func getEnv(key string, _default string) string {
	v := os.Getenv(key)
	if v == "" {
		return _default
	}
	return v
}
