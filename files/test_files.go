package files

import (
	"log"
	"testing"

	"github.com/dihedron/builds/files"
)

func TestExists(t *testing.T) {
	exists, err := files.Exists("main.go")
	if exists {
		log.Println("file exists")
	}
}
