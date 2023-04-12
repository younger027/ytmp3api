package covert

import (
	"regexp"
	"strings"
)

const (
	ErrInvalidCharactersInVideoID = constError("invalid characters in video id")
	ErrVideoIDMinLength           = constError("the video id must be at least 10 characters long")
)

type constError string

func (e constError) Error() string {
	return string(e)
}

var videoRegexpList = []*regexp.Regexp{
	regexp.MustCompile(`(?:v|embed|shorts|watch\?v)(?:=|/)([a-zA-Z0-9\-_]{11})`),
	regexp.MustCompile(`(?:=|/)([a-zA-Z0-9\-_]{11})`),
	regexp.MustCompile(`([a-zA-Z0-9\-_]{11})`),
}

// ExtractVideoID extracts the videoID from the given string
func ExtractVideoID(videoID string) (string, error) {
	if strings.Contains(videoID, "youtu") {
		for _, re := range videoRegexpList {
			if isMatch := re.MatchString(videoID); isMatch {
				subs := re.FindStringSubmatch(videoID)
				videoID = subs[1]
			}
		}
	}

	if strings.ContainsAny(videoID, "?&/<%=") {
		return "", ErrInvalidCharactersInVideoID
	}
	if len(videoID) < 10 {
		return "", ErrVideoIDMinLength
	}

	return videoID, nil
}
