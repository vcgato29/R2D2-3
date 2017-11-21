package main

import(
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var fileTypes = map[string]bool {
	".avi": true,
	".mp4": true,
	".mkv": true,
	".m4v": true,
	".wmv": true,
}


func main() {
	filepath.Walk("/mnt/media/downloads", processDirectory)
}

func processDirectory(path string, f os.FileInfo, err error) error {

	extension := filepath.Ext(f.Name())

	if fileTypes[extension] {
		name := strings.TrimSuffix(f.Name(), extension)
		content := MatchContent(name)

		fmt.Println(fmt.Sprintf("Name: %s, Episode: %s, Season: %s", content.name, content.number, content.season))
	}
	
	return nil
}