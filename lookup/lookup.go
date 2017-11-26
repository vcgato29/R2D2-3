package lookup

import(
	"github.com/ryanbradynd05/go-tmdb"
	"os"
	"errors"
)

type TvResult struct {
	Title string
	EpisodeName string
	EpisodeNum int
	SeasonNum int
	AirDate string	
	FirstAirDate string
}

func tmdbTvEpisode(seriesId int, season int, episode int) (*tmdb.TvEpisode, error) {
	db := tmdb.Init(os.Getenv("TMDB_API"))	
	lookup, err := db.GetTvEpisodeInfo(seriesId, season, episode, nil)

	if err == nil {
		return lookup, nil
	}

	dummy := tmdb.TvEpisode{}
	return &dummy, errors.New("no match at tmdb")
}

// Returns the most likely TV result from the TMdb API
func tmdbTv(show string) (*tmdb.TV, error) {
	db := tmdb.Init(os.Getenv("TMDB_API"))
	lookup, _ := db.SearchTv(show, nil)

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
		return db.GetTvInfo(candidateID, nil)
	}
	
	// Nothing found on TMdb
	return nil, errors.New("no TMdb match found when looking up show")
}

func Tv(show string, season int, episode int) (TvResult, error) {
	tmdbResult, err := tmdbTv(show)

	if err == nil {
		tmdbEpisode, err := tmdbTvEpisode(tmdbResult.ID, season, episode)

		if err == nil {			
			return TvResult{
				Title: tmdbResult.Name,
				EpisodeName: tmdbEpisode.Name,
				EpisodeNum: tmdbEpisode.EpisodeNumber,
				SeasonNum: tmdbEpisode.SeasonNumber,
				AirDate: tmdbEpisode.AirDate,
				FirstAirDate: tmdbResult.FirstAirDate,
				}, nil			
		}

		return TvResult{Title: "NA"}, err		

	}

	return TvResult{Title: "NA"}, err
}