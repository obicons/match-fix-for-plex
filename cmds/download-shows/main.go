package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/obicons/match-fix/showdb"
)

func main() {
	apiKey := flag.String("api-key", "", "Database key for https://www.themoviedb.org")
	flag.Parse()

	if *apiKey == "" {
		flag.Usage()
		os.Exit(1)
	}

	dbApi := showdb.NewMovieDB(*apiKey)
	shows, err := dbApi.DownloadShows()
	if err != nil {
		panic(err)
	}

	for _, show := range shows {
		fmt.Println(show)
	}

}
