package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/bnmcg/r2d2/lookup"
	"github.com/bnmcg/r2d2/matching"
)

// File types to organise
var fileTypes = map[string]bool{
	".avi": true,
	".mp4": true,
	".mkv": true,
	".m4v": true,
	".wmv": true,
}

// Common acronyms / abbreviations for films and TV shows that we can expand out to make
// matching easier later.
var acronyms = map[string]string{
	"himym":  "How I Met Your Mother",
	"taahm":  "Two and A Half Men",
	"tfpoba": "The Fresh Prince of Bel Air",
	"satc":   "Sex and The City",
	"tbbt":   "The Big Bang Theory",
	"potc":   "Pirates of The Caribbean",
	"rots":   "Return of The Sith",
	"aotc":   "Attack of The Clones",
}

var outputTvDirectoryFormat = "TV Shows/{Show}/Season {Season}/"
var outputTvFileFormat = "{Show} - S{Season}E{Number} - {Episode}"

var outputMovieDirectoryFormat = "Movies/{Movie} ({Year})/"
var outputMovieFileFormat = "{Movie} ({Year})"

var destinationDirectory string
var sourceDirectory string

func main() {
	// Get input directory
	// args[1] = Source path
	// args[2] = Destination path
	args := os.Args

	sourceDirectory = args[1]
	destinationDirectory = args[2]

	// Interval between runs
	interval, _ := strconv.Atoi(os.Getenv("R2D2_INTERVAL"))

	for {
		filepath.Walk(args[1], processDirectory)
		time.Sleep(time.Duration(interval) * time.Second)
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

		if content.Tv {
			result, err := lookup.Tv(content.Name, content.Season, content.Number)

			// Lookup successful
			if err == nil {
				fmt.Println(fmt.Sprintf("[tmdb] Show: %s, episode: %s, aired: %s, first seen: %s", result.Title, result.EpisodeName, result.AirDate, result.FirstAirDate))
				generatedFileName, generatedDirectory, generatedOutputPath := generateTvOutputPath(result, extension)
				fmt.Println(fmt.Sprintf("[r2d2] Generated output path: %s", generatedOutputPath))

				moveFile(path, generatedDirectory, generatedFileName)
			} else {
				fmt.Println(fmt.Sprintf("[r2d2] ERROR: %s", err))
				// Try swapping the matched episode title and series title
				result, err := lookup.Tv(content.Episode, content.Season, content.Number)
				if err == nil {
					fmt.Println(fmt.Sprintf("[tmdb] Show: %s, episode: %s, aired: %s, first seen: %s", result.Title, result.EpisodeName, result.AirDate, result.FirstAirDate))
					generatedFileName, generatedDirectory, generatedOutputPath := generateTvOutputPath(result, extension)
					fmt.Println(fmt.Sprintf("[r2d2] Generated output path: %s", generatedOutputPath))

					moveFile(path, generatedDirectory, generatedFileName)

				} else {
					fmt.Println(err)
				}
			}
		}

		if content.Movie {
			result, err := lookup.Movie(content.Name, content.Year)

			// Lookup successful
			if err == nil {
				fmt.Println(fmt.Sprintf("[tmdb] Movie: %s, released: %s", result.Title, result.ReleaseDate))
				generatedFileName, generatedDirectory, generatedOutputPath := generateMovieOutputPath(result, extension)
				fmt.Println(fmt.Sprintf("[r2d2] Generated output path: %s", generatedOutputPath))

				moveFile(path, generatedDirectory, generatedFileName)
			} else {
				fmt.Println(err)
			}
		}

	} else {
		// Remove files that aren't a valid file
		fmt.Println(fmt.Sprintf("[r2d2] Removing - not a valid extension: %s", path))
		if path != sourceDirectory {
			os.Remove(path)
		}
	}

	fmt.Println()
	time.Sleep(time.Second * 1)

	return nil
}

func moveFile(inputPath string, outputDirectoryPath string, outputFilename string) {
	// Move
	if !(os.Getenv("MODE") == "TEST") {
		fmt.Println("[r2d2] Moving files")
		// Create required directory
		os.MkdirAll(outputDirectoryPath, os.ModePerm)
		commandArgs := []string{inputPath, fmt.Sprintf("%s%s", outputDirectoryPath, outputFilename)}
		moveError := exec.Command("mv", commandArgs...).Run()

		if moveError != nil {
			fmt.Println(moveError)
		}
	}
}

func generateTvOutputPath(lookup lookup.TvResult, extension string) (string, string, string) {
	outputPath := strings.Replace(outputTvDirectoryFormat, "{Show}", lookup.Title, -1)
	outputPath = strings.Replace(outputPath, "{Season}", fmt.Sprintf("%02d", lookup.SeasonNum), -1)

	outputFile := strings.Replace(outputTvFileFormat, "{Show}", lookup.Title, -1)
	outputFile = strings.Replace(outputFile, "{Episode}", lookup.EpisodeName, -1)
	outputFile = strings.Replace(outputFile, "{Season}", fmt.Sprintf("%02d", lookup.SeasonNum), -1)
	outputFile = strings.Replace(outputFile, "{Number}", fmt.Sprintf("%02d", lookup.EpisodeNum), -1)
	outputFile += extension

	// Return file, directory, combined
	return outputFile, fmt.Sprintf("%s%s", destinationDirectory, outputPath), fmt.Sprintf("%s%s%s", destinationDirectory, outputPath, outputFile)
}

func generateMovieOutputPath(lookup lookup.MovieResult, extension string) (string, string, string) {
	outputPath := strings.Replace(outputMovieDirectoryFormat, "{Movie}", lookup.Title, -1)
	outputPath = strings.Replace(outputPath, "{Year}", strconv.Itoa(lookup.Year), -1)

	outputFile := strings.Replace(outputMovieFileFormat, "{Movie}", lookup.Title, -1)
	outputFile = strings.Replace(outputFile, "{Year}", strconv.Itoa(lookup.Year), -1)

	outputFile += extension

	// Return file, directory, combined
	return outputFile, fmt.Sprintf("%s%s", destinationDirectory, outputPath), fmt.Sprintf("%s%s%s", destinationDirectory, outputPath, outputFile)
}
