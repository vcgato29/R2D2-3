package lookup

import(
	"github.com/ryanbradynd05/go-tmdb"
	"os"
	"errors"
)

type Result struct {
	Title string
	ReleaseDate string
	Genres []struct {ID int; Name string }
}

// Returns the most likely TV result from the TMdb API
func tmdbMovie(name string) (*tmdb.Movie, error) {
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

func Movie(name string, year int) (Result, error) {
	tmdbResult, err := tmdbMovie(name)

	if err == nil {
		return Result{
			Title: tmdbResult.Title,
			ReleaseDate: tmdbResult.ReleaseDate,
			Genres: tmdbResult.Genres,
		}, err		

	}

	return Result{Title: "NA"}, err
}