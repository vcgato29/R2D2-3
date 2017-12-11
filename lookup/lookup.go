package lookup

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ferhatelmas/levenshtein"
	"github.com/ryanbradynd05/go-tmdb"
)

type TvResult struct {
	Title        string
	EpisodeName  string
	EpisodeNum   int
	SeasonNum    int
	AirDate      string
	FirstAirDate string
}

func tmdbTvEpisode(seriesID int, season int, episode int) (*tmdb.TvEpisode, error) {
	db := tmdb.Init(os.Getenv("TMDB_API"))
	lookup, err := db.GetTvEpisodeInfo(seriesID, season, episode, nil)

	if err == nil {
		return lookup, nil
	}

	dummy := tmdb.TvEpisode{}
	return &dummy, errors.New("no match at tmdb")
}

// Returns the most likely TV result from the TMdb API
func tmdbTv(show string, episodeNumber int, seasonNumber int) (*tmdb.TV, error) {
	db := tmdb.Init(os.Getenv("TMDB_API"))
	show = strings.Replace(show, ".", " ", -1)
	fmt.Println(fmt.Sprintf("[r2d2] Looking up show: %s", show))
	lookup, _ := db.SearchTv(show, nil)

	candidateID := 0
	candidateDist := 1000

	for _, element := range lookup.Results {
		if candidateID == 0 {
			fmt.Println(fmt.Sprintf("Initial show %s (%s) / %d", element.Name, element.FirstAirDate, element.ID))

			candidateID = element.ID
			candidateDist = levenshtein.Dist(show, element.Name)
		} else {
			// More popular shows are more likely to be added to media libarires.
			// Particularly if more popular with more votes
			fmt.Println(fmt.Sprintf("Found show %s (%s) / %d", element.Name, element.FirstAirDate, element.ID))
			fmt.Println(fmt.Sprintf("Considering show %s (%s) / %d", element.Name, element.FirstAirDate, element.ID))
			// Check show particulars
			info, _ := db.GetTvInfo(candidateID, nil)

			for _, season := range info.Seasons {
				if season.SeasonNumber != seasonNumber {
					continue
				}

				if season.EpisodeCount >= episodeNumber {
					fmt.Println(fmt.Sprintf("[r2d2] Inspecting season %d with total number of episodes %d", season.SeasonNumber, season.EpisodeCount))
					if levenshtein.Dist(show, element.Name) < candidateDist {
						candidateID = element.ID
					}
				}
			}
		}
	}

	// Provided that one result was found, get a full listing for that element
	if candidateID != 0 {
		fmt.Println(fmt.Sprintf("[tmdb] Final show ID %d", candidateID))
		return db.GetTvInfo(candidateID, nil)
	}

	// Nothing found on TMdb
	return nil, errors.New("no TMdb match found when looking up show")
}

func Tv(show string, season int, episode int) (TvResult, error) {
	fmt.Println(fmt.Sprintf("[r2d2] Looking up %s, season %d, episode %d", show, season, episode))
	tmdbResult, err := tmdbTv(show, episode, season)

	if err == nil {
		tmdbEpisode, err := tmdbTvEpisode(tmdbResult.ID, season, episode)

		if err == nil {
			return TvResult{
				Title:        tmdbResult.Name,
				EpisodeName:  tmdbEpisode.Name,
				EpisodeNum:   tmdbEpisode.EpisodeNumber,
				SeasonNum:    tmdbEpisode.SeasonNumber,
				AirDate:      tmdbEpisode.AirDate,
				FirstAirDate: tmdbResult.FirstAirDate,
			}, nil
		}

		return TvResult{Title: "NA"}, err

	}

	fmt.Println(fmt.Sprintf("[r2d2] ERROR: %s", err))
	return TvResult{Title: "NA"}, err
}
