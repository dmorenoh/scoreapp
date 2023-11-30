package score

import (
	"errors"
	"sort"
)

type Entity struct {
	User    uint `json:"user"`
	Total   int  `json:"total"`
	Version uint `json:"version"`
}

type entitiesMap map[uint]*Entity

func (e entitiesMap) sortValues() Entities {
	scoreList := make([]*Entity, 0, len(e))
	for _, score := range e {
		scoreList = append(scoreList, score)
	}

	sort.Slice(scoreList, func(i, j int) bool {
		return scoreList[i].Total > scoreList[j].Total
	})
	return scoreList

}

type Entities []*Entity

func (e Entities) filter(filter interface{}) (Entities, error) {
	switch f := filter.(type) {
	case Absolute:
		if f.Limit > uint(len(e)) {
			return nil, errors.New("limit is greater than the number of scores")
		}
		return e[0:f.Limit], nil
	case Relative:
		if f.Position-f.Around < 0 || f.Position+f.Around >= uint(len(e)) {
			return nil, errors.New("position is greater than the number of scores")
		}
		return e[f.Position-f.Around : f.Position+f.Around+1], nil
	default:
		return nil, errors.New("invalid filter")
	}
}
