//go:build !prod
// +build !prod

package web

import (
	"fmt"
	"io/fs"
	"os"
)

func GetFileSystem() (fs.FS, error) {
	fmt.Println("Filesystem served by os package")
	return os.DirFS("web/dist"), nil
}
