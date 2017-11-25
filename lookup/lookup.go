package lookup

import(
	"github.com/ryanbradynd05/go-tmdb"
	"os"
	"fmt"
)

func LookupTmdbTv(query string) {
	db := tmdb.Init(os.Getenv("TMDB_API"))
	results, _ := db.SearchTv(query, nil)

	fmt.Println(results.Results[0].Name)
	fmt.Println(results.Results[0].Popularity)
	fmt.Println(results.Results[0].VoteAverage)
}

func LookupTmdbMovie(query string) {
	db := tmdb.Init(os.Getenv("TMDB_API"))
	results, _ := db.SearchTv(query, nil)

	fmt.Println(results.Results[0].Name)
	fmt.Println(results.Results[0].Popularity)
	fmt.Println(results.Results[0].VoteAverage)
}