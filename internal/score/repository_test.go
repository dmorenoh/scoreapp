package score

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestInMemoryRepository_Save_Race(t *testing.T) {
	scores := map[uint]*Entity{}
	repo := NewInMemoryRepository(scores)

	var wg sync.WaitGroup

	var errors int
	var start = make(chan bool)
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			<-start
			score := NewScore(123, i*10)
			err := repo.Save(nil, score)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				errors++
			}
			wg.Done()
		}()
	}

	close(start)
	wg.Wait()

	assert.Equal(t, 2, errors)
}
