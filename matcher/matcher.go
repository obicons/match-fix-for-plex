package matcher

import (
	"bufio"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type emptyStruct struct{}

var fileTypes = []string{
	"asf",
	"mp3",
	"mp4",
	"mkv",
	"wav",
	"alac",
	"aac",
	"avi",
}

var stopWords = map[string]interface{}{
	"brrip":    emptyStruct{},
	"complete": emptyStruct{},
	"episode":  emptyStruct{},
	"season":   emptyStruct{},
	"web":      emptyStruct{},
	"xvid":     emptyStruct{},
}

/*
 * Returns the name of the TV show that best matches directoryName
 */
func MatchTV(directoryName string, br *bufio.Reader) string {
	directoryName = strings.ToLower(directoryName)
	directoryName = removeSeasonDataOn(directoryName)
	directoryName = removeLikelyMetaData(directoryName)
	_, directoryName = removeMediaExtension(directoryName)
	directoryName = replaceWithSpaces(directoryName)
	directoryName = removeStopWords(directoryName)
	directoryName = strings.TrimSpace(directoryName)

	bestMatch := ""
	bestScore := -1
	scanner := bufio.NewScanner(br)
	for scanner.Scan() {
		title := strings.ToLower(scanner.Text())
		score := lcs(title, directoryName)
		if score > bestScore {
			bestScore = score
			bestMatch = title
		}
	}

	return bestMatch
}

/*
 * Extracts just the season component from a title
 * Returns -1 if no season number can be extracted
 */
func MatchSeasonNo(title string) int {
	title = strings.ToLower(title)
	title = replaceWithSpaces(title)

	r := regexp.MustCompile("season ?([0-9]+)")
	matches := r.FindStringSubmatch(title)
	if len(matches) == 2 {
		if seasonNo, err := strconv.Atoi(matches[1]); err == nil {
			return seasonNo
		}
	}

	return -1
}

/*
 * Removes data within [] and () from string
 */
func removeLikelyMetaData(name string) string {
	r := regexp.MustCompile("(\\[.*\\]|\\(.*\\))")
	return r.ReplaceAllString(name, "")
}

/*
 * Removes season information
 */
func removeSeasonDataOn(name string) string {
	r := regexp.MustCompile("season.*")
	return r.ReplaceAllString(name, "")
}

/*
 * Removes media file extensions
 */
func removeMediaExtension(filename string) (extension, basename string) {
	for _, ext := range fileTypes {
		if strings.HasSuffix(filename, "."+ext) {
			extension = ext
			basename = strings.TrimSuffix(filename, "."+ext)
			return
		}
	}
	basename = filename
	return
}

/*
 * Removes stop-words from the name
 */
func removeStopWords(name string) string {
	var builder strings.Builder
	var sep byte
	words := strings.Split(name, " ")
	for _, word := range words {
		if word == "" {
			continue
		}
		if _, isStop := stopWords[word]; !isStop {
			if sep != 0 {
				builder.WriteByte(sep)
			}
			builder.WriteString(word)
		}
		sep = ' '
	}
	return builder.String()
}

/*
 * Replaces non-letters with spaces
 */
func replaceWithSpaces(name string) string {
	var builder strings.Builder
	for _, char := range name {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			builder.WriteByte(' ')
		} else {
			builder.WriteRune(char)
		}
	}
	return builder.String()
}

/*
 * Returns if a[i] is the start of a word
 */
func atWordStart(a string, i int) bool {
	if i == 0 {
		return true
	}

	return a[i-1] == ' '
}

/*
 * Scores the quality of a match
 */
func scoreMatch(a, b string, i, j int) int {
	if atWordStart(a, i) {
		// reward start-of-word matches more
		return 2
	} else if i > 0 && j > 0 && a[i-1] == b[j-1] {
		// reward consecutive matches more
		return 2
	}
	return 1
}

/*
 * Returns the length of the longest common subsequence of a and b
 */
func lcs(a, b string) int {
	// garbage in, garbage out
	if len(a) == 0 || len(b) == 0 {
		return 0
	}

	subsequenceLengths := make([][]int, len(a)+1)
	for i := range subsequenceLengths {
		subsequenceLengths[i] = make([]int, len(b)+1)
	}

	// initial conditions: i or j = 0 => lcs(i, j) = 0
	for j := 0; j < len(b); j++ {
		subsequenceLengths[0][j] = 0
	}
	for i := 0; i < len(a); i++ {
		subsequenceLengths[i][0] = 0
	}

	// main logic: compute lcs matrix
	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			if a[i-1] == b[j-1] {
				subsequenceLengths[i][j] = scoreMatch(a, b, i-1, j-1) +
					subsequenceLengths[i-1][j-1]
			} else {
				subsequenceLengths[i][j] = max(
					subsequenceLengths[i-1][j],
					subsequenceLengths[i][j-1],
				)
			}
		}
	}

	return subsequenceLengths[len(a)][len(b)]
}

/*
 * Returns the maximum of two ints
 */
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
