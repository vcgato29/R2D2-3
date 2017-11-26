package main

import(
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"github.com/bnmcg/r2d2/lookup"
	"github.com/bnmcg/r2d2/matching"
	"time"
)

// File types to organise
var fileTypes = map[string]bool {
	".avi": true,
	".mp4": true,
	".mkv": true,
	".m4v": true,
	".wmv": true,
}

// Common acronyms / abbreviations for films and TV shows that we can expand out to make
// matching easier later.
var acronyms = map[string]string {
	"himym": "How I Met Your Mother",
	"taahm": "Two and A Half Men",
	"tfpoba": "The Fresh Prince of Bel Air",
	"satc": "Sex and The City",
	"tbbt": "The Big Bang Theory",
	"potc": "Pirates of The Caribbean",
	"rots": "Return of The Sith",
	"aotc": "Attack of The Clones",
}

var outputDirectoryFormat = "{Show}/Season {Season}/"
var outputFileFormat = "{Show} - S{Season}E{Number} - {Episode}"

func main() {
	// Get input directory
	// args[1] = Source path
	args := os.Args
	for {
		filepath.Walk(args[1], processDirectory)
		time.Sleep(time.Second * 30)	
	}
}

func processDirectory(path string, f os.FileInfo, err error) error {

	fmt.Println(fmt.Sprintf("[r2d2] Processing file: %s", path))

	extension := filepath.Ext(f.Name())
	fmt.Println(fmt.Sprintf("[r2d2] Matched extension: %s", extension))	

	if fileTypes[extension] {
		name := strings.TrimSuffix(f.Name(), extension)
		fmt.Println(fmt.Sprintf("[r2d2] File name: %s", name))			

		content := matching.MatchContent(name)

		result, err := lookup.Tv(content.Name, content.Season, content.Number)
		
		// Lookup successful
		if err == nil {
			fmt.Println(fmt.Sprintf("[tmdb] Show: %s, episode: %s, aired: %s, first seen: %s", result.Title, result.EpisodeName, result.AirDate, result.FirstAirDate))
			// Generate output path
			outputPath := strings.Replace(outputDirectoryFormat, "{Show}", result.Title, -1)
			outputPath = strings.Replace(outputPath, "{Season}", fmt.Sprintf("%02d", result.SeasonNum), -1)

			outputFile := strings.Replace(outputFileFormat, "{Show}", result.Title, -1)
			outputFile = strings.Replace(outputFile, "{Episode}", result.EpisodeName, -1)			
			outputFile = strings.Replace(outputFile, "{Season}", fmt.Sprintf("%02d", result.SeasonNum), -1)
			outputFile = strings.Replace(outputFile, "{Number}", fmt.Sprintf("%02d", result.EpisodeNum), -1)		
			outputFile += extension

			fmt.Println(fmt.Sprintf("[r2d2] Generated output path: %s%s", outputPath, outputFile))
		// Lookup failed
		} else {
			// Try swapping the matched episode title and series title
			result, err := lookup.Tv(content.Episode, content.Season, content.Number)
			if err == nil {
				fmt.Println(fmt.Sprintf("[tmdb] Show: %s, episode: %s, aired: %s, first seen: %s", result.Title, result.EpisodeName, result.AirDate, result.FirstAirDate))
			// Give up
			} else {
				fmt.Println(err)		
			}			
		}
	} else {
		fmt.Println(fmt.Sprintf("[r2d2] Skipping - not a valid extension: %s", path))
	}

	fmt.Println()
	time.Sleep(time.Second * 1)
	
	return nil
}

func generateTvOutputPath(lookup lookup.TvResult) string {
	return ""
}