package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/obicons/match-fix/matcher"
)

func main() {
	titlePath := flag.String("titles", "", "path to file containing TV Show titles")
	toMatch := flag.String("matchee", "", "title that needs matched")
	flag.Parse()

	if *titlePath == "" || *toMatch == "" {
		flag.Usage()
		os.Exit(1)
	}

	if _, err := os.Stat(*toMatch); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	// matchDir - the name of the TV show target we're going to match
	matchDir := path.Base(*toMatch)
	// match Loc - the name of the location on the filesystem of the target
	matchLoc := path.Dir(*toMatch)

	file, err := os.Open(*titlePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	br := bufio.NewReader(file)
	match := matcher.MatchTV(matchDir, br)

	// check if the match is good with the user
	shouldRename := PromptYorN("Match \"%s\" with \"%s\"?", *toMatch, match)
	if !shouldRename {
		fmt.Println("No suitable match found. Exiting.")
		os.Exit(1)
	}

	// find the season #
	var newLoc string
	seasonNo := matcher.MatchSeasonNo(matchDir)
	if seasonNo == -1 {
		newLoc = path.Join(matchLoc, match)
	} else {
		newLoc = path.Join(matchLoc, match, fmt.Sprintf("Season %d", seasonNo))
	}

	// check that the move is good with the user
	shouldMove := PromptYorN("Move %s to %s?", *toMatch, newLoc)
	if !shouldMove {
		fmt.Println("No more repairs possible. Exiting.")
		os.Exit(1)
	}

	// verify the output directory already exists.
	// create it, if it doesn't
	moveDir := path.Dir(newLoc)
	if _, err := os.Stat(moveDir); os.IsNotExist(err) {
		shouldCreate := PromptYorN("%s does not exist. Create it?", moveDir)
		if !shouldCreate {
			fmt.Println("No more repairs possible. Exiting.")
			os.Exit(1)
		}
		os.Mkdir(moveDir, 0777)
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}

	err = os.Rename(*toMatch, newLoc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: moving failed (%s). Verify your directory integrity.\n", err)
		os.Exit(1)
	}
}

/*
 * Repeatedly prompts a yes or no question by passing msg and args to printf.
 * Returns true if the user says yes, otherwise false.
 */
func PromptYorN(msg string, args ...interface{}) bool {
	hadValidResponse := false
	resp := false
	br := bufio.NewScanner(os.Stdin)

	for !hadValidResponse {
		fmt.Printf(msg+" [Y/N]: ", args...)
		if !br.Scan() {
			break
		}
		response := br.Text()

		if response == "y" || response == "yes" {
			hadValidResponse = true
			resp = true
		} else if response == "n" || response == "no" {
			hadValidResponse = true
		}
	}

	return resp
}
