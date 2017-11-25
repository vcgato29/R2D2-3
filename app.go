package main

import(
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"github.com/bnmcg/r2d2/lookup"
	"github.com/bnmcg/r2d2/matching"
)

var fileTypes = map[string]bool {
	".avi": true,
	".mp4": true,
	".mkv": true,
	".m4v": true,
	".wmv": true,
}


func main() {
	//filepath.Walk("/mnt/media/downloads", processDirectory)
	lookup.LookupTmdbTv("Lif")
}

func processDirectory(path string, f os.FileInfo, err error) error {

	extension := filepath.Ext(f.Name())

	if fileTypes[extension] {
		name := strings.TrimSuffix(f.Name(), extension)
		content := matching.MatchContent(name)

		fmt.Println(fmt.Sprintf("Name: %s, Episode: %s, Season: %s", content.Name, content.Number, content.Season))
	}
	
	return nil
}