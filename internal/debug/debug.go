package debug

import (
	"log"
	"os"
)

// Debug fonction
func Debug(format string, v ...any) {
	debug := os.Getenv("DEBUG")

	if debug == "1" {
		log.Printf(format, v...)
	}
}
