package score

import (
	"context"
)

type Service interface {
	SubmitAbsolute(ctx context.Context, user uint, total int) error
	SubmitRelative(ctx context.Context, user uint, variation int) error
	Find(ctx context.Context, filter interface{}) ([]View, error)
}

type service struct {
	repo Repository
}

func NewService(repository Repository) Service {
	return &service{
		repo: repository,
	}
}

func (s *service) SubmitAbsolute(ctx context.Context, user uint, total int) error {
	scr, err := s.repo.Get(ctx, user)
	if err != nil {
		return err
	}

	if scr == nil {
		scr = NewScore(user, total)
	} else {
		scr.Total(total)
	}

	if err := s.repo.Save(ctx, scr); err != nil {
		return err
	}
	return nil
}

func (s *service) SubmitRelative(ctx context.Context, user uint, variation int) error {
	scr, err := s.repo.Get(ctx, user)
	if err != nil {
		return err
	}

	if scr == nil {
		scr = NewScore(user, variation)
	} else {
		scr.AddScore(variation)
	}

	if err := s.repo.Save(ctx, scr); err != nil {
		return err
	}
	return nil
}

func (s *service) Find(ctx context.Context, filter interface{}) ([]View, error) {
	scores, err := s.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	response := make([]View, len(scores))
	for i, score := range scores {
		response[i] = NewView(score)
	}

	return response, err
}

type Absolute struct {
	Limit uint
}

type Relative struct {
	Position uint
	Around   uint
}
