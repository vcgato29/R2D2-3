package lookup

import (
	"errors"
	"os"
	"strings"

	"github.com/araddon/dateparse"
	"github.com/ryanbradynd05/go-tmdb"
)

type MovieResult struct {
	Title       string
	ReleaseDate string
	Year        int
	Genres      []struct {
		ID   int
		Name string
	}
}

// Returns the most likely TV result from the TMdb API
func tmdbMovie(name string) (*tmdb.Movie, error) {
	// Replace any .'s in the title with spaces
	name = strings.Replace(name, ".", " ", -1)

	db := tmdb.Init(os.Getenv("TMDB_API"))
	lookup, _ := db.SearchMovie(name, nil)

	candidateID := 0
	var candidateRating float32
	var candidateVotes uint32

	for _, element := range lookup.Results {
		if candidateID == 0 {
			candidateID = element.ID
			candidateRating = element.VoteAverage
			candidateVotes = element.VoteCount
		} else {
			// More popular shows are more likely to be added to media libarires.
			// Particularly if more popular with more votes
			if element.VoteAverage > candidateRating && element.VoteCount > candidateVotes {
				candidateID = element.ID
				candidateRating = element.VoteAverage
				candidateVotes = element.VoteCount
			}
		}
	}

	// Provided that one result was found, get a full listing for that element
	if candidateID != 0 {
		return db.GetMovieInfo(candidateID, nil)
	}

	// Nothing found on TMdb
	return nil, errors.New("no TMdb match found when looking up movie")
}

func Movie(name string, year int) (MovieResult, error) {
	tmdbResult, err := tmdbMovie(name)

	if err == nil {

		// Parse release date
		date, parseError := dateparse.ParseAny(tmdbResult.ReleaseDate)

		if parseError == nil {
			return MovieResult{
				Title:       tmdbResult.Title,
				ReleaseDate: tmdbResult.ReleaseDate,
				Year:        date.Year(),
				Genres:      tmdbResult.Genres,
			}, err
		}

		return MovieResult{
			Title:       tmdbResult.Title,
			ReleaseDate: tmdbResult.ReleaseDate,
			Year:        0000,
			Genres:      tmdbResult.Genres,
		}, err
	}

	return MovieResult{Title: "NA"}, err
}
