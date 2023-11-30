package server

import (
	"errors"
	"regexp"
	"scoreapp/internal/score"
	"strconv"
	"strings"
)

func BuildFilter(value string) (interface{}, error) {
	abs, _ := regexp.Compile("top\\d*")
	rel, _ := regexp.Compile("At\\d*/\\d")

	switch {
	case abs.MatchString(value):
		limit, err := strconv.Atoi(value[3:])
		if err != nil {
			return nil, err
		}
		return score.Absolute{
			Limit: uint(limit),
		}, nil
	case rel.MatchString(value):
		segments := strings.Split(value, "/")
		position, err := strconv.Atoi(segments[0][2:])
		if err != nil {
			return nil, err
		}
		around, err := strconv.Atoi(segments[1])
		return score.Relative{
			Position: uint(position),
			Around:   uint(around),
		}, nil
	}
	return nil, errors.New("invalid filter")
}

func ScoreVariationValue(value string) (int, error) {
	r, _ := regexp.Compile("[+|\\-]\\d*")
	if !r.MatchString(value) {
		return 0, errors.New("invalid score variation")
	}
	sign := value[0]
	number, err := strconv.Atoi(value[1:])
	if err != nil {
		return 0, err
	}

	if sign == '-' {
		return -number, nil
	}

	return number, nil

}
