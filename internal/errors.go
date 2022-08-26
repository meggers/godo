package internal

import (
	"log"
)

func CheckError(err error, message string) {
	if err != nil {
		log.Fatalf("%v with: %v\n", message, err)
	}
}
