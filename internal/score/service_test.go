package score

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestService_SubmitAbsolute(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	existingScore := Entity{User: 123, Total: 100, Version: 1}
	scores := map[uint]*Entity{existingScore.User: &existingScore}
	inMemRepo := &inMemoryRepository{
		scores: scores,
	}

	type args struct {
		user  uint
		total int
	}

	tests := []struct {
		name  string
		repo  Repository
		args  args
		error error
	}{
		{
			name: "submit absolute score for an existing user",
			repo: inMemRepo,
			args: args{
				user:  123,
				total: 100,
			},
		},
		{
			name: "submit absolute score for an non-existing user",
			repo: inMemRepo,
			args: args{
				user:  345,
				total: 100,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				repo: tt.repo,
			}
			err := s.SubmitAbsolute(ctx, tt.args.user, tt.args.total)
			assert.ErrorIs(t, err, tt.error)
			f := scores[tt.args.user]
			assert.Equal(t, tt.args.total, f.Total)
		})
	}
}

func TestService_SubmitRelative(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	existingScore := Entity{User: 123, Total: 100, Version: 1}
	scores := map[uint]*Entity{existingScore.User: &existingScore}
	inMemRepo := &inMemoryRepository{
		scores: scores,
	}

	type args struct {
		user      uint
		variation int
	}

	tests := []struct {
		name          string
		repo          Repository
		args          args
		expectedTotal int
		error         error
	}{
		{
			name: "submit relative score for an existing user",
			repo: inMemRepo,
			args: args{
				user:      123,
				variation: 100,
			},
			expectedTotal: 200,
		},
		{
			name: "submit relative score for an non-existing user",
			repo: inMemRepo,
			args: args{
				user:      345,
				variation: 100,
			},
			expectedTotal: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				repo: tt.repo,
			}
			err := s.SubmitRelative(ctx, tt.args.user, tt.args.variation)
			assert.ErrorIs(t, err, tt.error)
			f := scores[tt.args.user]
			assert.Equal(t, tt.expectedTotal, f.Total)
		})
	}
}

func TestService_Find(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	scores := make(map[uint]*Entity)

	for i := 0; i < 100; i++ {
		scores[uint(i)] = &Entity{
			User:  uint(i + 1),
			Total: (i + 1) * 10,
		}
	}

	inMemRepo := &inMemoryRepository{
		scores: scores,
	}

	type args struct {
		filter interface{}
	}

	tests := []struct {
		name           string
		repo           Repository
		args           args
		error          error
		expectedLength int
	}{
		{
			name: "find user scores by absolute filter",
			repo: inMemRepo,
			args: args{
				filter: Absolute{
					Limit: 50,
				},
			},
			expectedLength: 50,
		},
		{
			name: "find user scores by relative filter",
			repo: inMemRepo,
			args: args{
				filter: Relative{
					Position: 20,
					Around:   3,
				},
			},
			expectedLength: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				repo: tt.repo,
			}
			result, err := s.Find(ctx, tt.args.filter)
			assert.ErrorIs(t, err, tt.error)
			assert.NotNil(t, result)
			assert.Len(t, result, tt.expectedLength)
		})
	}
}
