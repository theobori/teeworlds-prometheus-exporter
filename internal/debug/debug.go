package debug

import (
	"log"
	"os"
)

func Debug(format string, v ...any) {
	debug := os.Getenv("DEBUG")

	if debug == "1" {
		log.Printf(format, v...)
	}
}
