package matcher

import "testing"

func TestMatchSeasonNoWithSpaces(t *testing.T) {
	title := "Hello World Season 3"
	seasonNo := MatchSeasonNo(title)
	expectedSeasonNo := 3
	if seasonNo != expectedSeasonNo {
		t.Fatalf("Expected season %d, found %d", expectedSeasonNo, seasonNo)
	}
}

func TestMatchSeasonNoWithDots(t *testing.T) {
	title := "Hello.World.Season.3"
	seasonNo := MatchSeasonNo(title)
	expectedSeasonNo := 3
	if seasonNo != expectedSeasonNo {
		t.Fatalf("Expected season %d, found %d", expectedSeasonNo, seasonNo)
	}
}

func TestMatchSeasonTitleWithNoSeason(t *testing.T) {
	title := "Hello.World.Season"
	seasonNo := MatchSeasonNo(title)
	expectedSeasonNo := -1
	if seasonNo != expectedSeasonNo {
		t.Fatalf("Expected season %d, found %d", expectedSeasonNo, seasonNo)
	}
}

func TestMatchSeasonTitleWithInvalidSeason(t *testing.T) {
	title := "Hello.World.Season this test is good"
	seasonNo := MatchSeasonNo(title)
	expectedSeasonNo := -1
	if seasonNo != expectedSeasonNo {
		t.Fatalf("Expected season %d, found %d", expectedSeasonNo, seasonNo)
	}
}

func TestMatchSeasonTitleWithSeasonNoSpace(t *testing.T) {
	title := "Hello.World.Season3"
	seasonNo := MatchSeasonNo(title)
	expectedSeasonNo := 3
	if seasonNo != expectedSeasonNo {
		t.Fatalf("Expected season %d, found %d", expectedSeasonNo, seasonNo)
	}
}
