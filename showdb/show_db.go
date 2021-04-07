package showdb

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// interface for interacting with a database storing TV show titles
type ShowDB interface {
	// returns an array of show titles and possibly an error
	DownloadShows() ([]string, error)
}

// interacts with https://www.themoviedb.org/
type MovieDB struct {
	// store the developer's API key
	apiKey string
}

// stores the result of an API call to MovieDB
type MovieDBResponse struct {
	Results []struct {
		Name string `json:"name"`
	} `json:"results"`
}

// generic authorization error
var AuthError = errors.New("Invalid API Authorization")

// URL of the movie database's API
const movieDBURL = `https://api.themoviedb.org/3/`

// limit to page number
const movieDBURLPageLimit = 500

// create a new MovieDB instance
func NewMovieDB(apiKey string) *MovieDB {
	return &MovieDB{apiKey: apiKey}
}

// implements showDB
func (db *MovieDB) DownloadShows() ([]string, error) {
	var shows []string
	for i := 0; i < movieDBURLPageLimit; i++ {
		round, err := db.downloadPageNo(i)
		if err != nil {
			return []string{}, err
		}
		shows = append(shows, round...)
	}
	return shows, nil
}

/*
 * Downloads the specified page number.
 * If the pageNo is out of range, returns an empty array.
 * If the apiKey is invalid, returns an appropriate error.
 * If the authorization is invalid, returns AuthError.
 */
func (db *MovieDB) downloadPageNo(pageNo int) ([]string, error) {
	// garbage in, garbage out
	if pageNo > movieDBURLPageLimit {
		return []string{}, nil
	}

	// build the request
	apiURLPath := movieDBURL + `discover/tv`
	apiURL, err := url.Parse(apiURLPath)
	AssertNever(err)

	params := make(url.Values)
	params.Add("page", strconv.Itoa(pageNo))
	params.Add("api_key", db.apiKey)
	apiURL.RawQuery = params.Encode()

	// make the request
	resp, err := http.Get(apiURL.String())
	if err != nil {
		return []string{}, err
	} else if resp.StatusCode == http.StatusUnauthorized {
		return []string{}, AuthError
	} else if resp.StatusCode == 404 {
		// TODO: does an out of range page 404?
		return []string{}, nil
	}

	// decode the response
	var titles []string
	var apiResp MovieDBResponse
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&apiResp)
	for _, show := range apiResp.Results {
		titles = append(titles, show.Name)
	}

	return titles, nil
}

func AssertNever(err error) {
	if err != nil {
		panic(err)
	}
}
