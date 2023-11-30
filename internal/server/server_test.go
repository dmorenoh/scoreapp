package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"scoreapp/internal/score"
	"testing"
)

func TestNewHTTPServer_SubmitScore(t *testing.T) {
	tests := []struct {
		name                 string
		request              string
		bodyRequest          string
		service              score.Service
		expectedHttpCode     int
		expectedBodyResponse string
	}{
		{
			name:                 "submit absolute score",
			request:              "POST",
			bodyRequest:          `{"user": 123, "total": 100}`,
			service:              &fakeOKService{},
			expectedHttpCode:     http.StatusCreated,
			expectedBodyResponse: `{"message":"accepted"}`,
		},
		{
			name:                 "submit relative score",
			request:              "POST",
			bodyRequest:          `{"user": 345, "score": "+100"}`,
			service:              &fakeOKService{},
			expectedHttpCode:     http.StatusCreated,
			expectedBodyResponse: `{"message":"accepted"}`,
		},
		{
			name:                 "submit relative score",
			request:              "POST",
			bodyRequest:          `{"user": 345, "score": "-100"}`,
			service:              &fakeOKService{},
			expectedHttpCode:     http.StatusCreated,
			expectedBodyResponse: `{"message":"accepted"}`,
		},
		{
			name:                 "submit relative score invalid request",
			request:              "POST",
			bodyRequest:          `{"user": 345, "score": "100"}`,
			service:              &fakeOKService{},
			expectedHttpCode:     http.StatusBadRequest,
			expectedBodyResponse: `{"message":"invalid json"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := NewApp(tt.service)
			router := SetupRouter(app)

			req, _ := http.NewRequest("POST", "/user/123/score", bytes.NewBuffer([]byte(tt.bodyRequest)))
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)
			assert.Equal(t, tt.expectedHttpCode, recorder.Code)
			assert.Equal(t, tt.expectedBodyResponse, recorder.Body.String())
		})
	}
}

func TestNewHTTPServer_GetRanking(t *testing.T) {
	scores := []score.View{
		{User: 1, Total: 100},
		{User: 2, Total: 100},
		{User: 3, Total: 100},
		{User: 4, Total: 100},
		{User: 5, Total: 100},
	}
	tests := []struct {
		name                 string
		request              string
		bodyRequest          string
		service              score.Service
		expectedHttpCode     int
		expectedBodyResponse []score.View
	}{
		{
			name:                 "get ranking absolute request",
			request:              "GET",
			bodyRequest:          `type=top10`,
			service:              &fakeOKService{scores: scores},
			expectedHttpCode:     http.StatusOK,
			expectedBodyResponse: scores,
		},
		{
			name:                 "get ranking relative request",
			request:              "GET",
			bodyRequest:          `type=At100/3`,
			service:              &fakeOKService{scores: scores},
			expectedHttpCode:     http.StatusOK,
			expectedBodyResponse: scores,
		},
		{
			name:             "get ranking invalid absolute request",
			request:          "GET",
			bodyRequest:      `type=topp10`,
			service:          &fakeOKService{scores: scores},
			expectedHttpCode: http.StatusBadRequest,
		},
		{
			name:             "get ranking invalid relative request",
			request:          "GET",
			bodyRequest:      `type=Atp100/3`,
			service:          &fakeOKService{scores: scores},
			expectedHttpCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := NewApp(tt.service)
			router := SetupRouter(app)

			url := fmt.Sprintf("/ranking?%s", tt.bodyRequest)
			req, _ := http.NewRequest("GET", url, nil)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)
			assert.Equal(t, tt.expectedHttpCode, recorder.Code)

			if tt.expectedHttpCode == http.StatusOK {
				var results []score.View
				body := recorder.Body
				err := json.Unmarshal(body.Bytes(), &results)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBodyResponse, results)
			}
		})
	}
}

type fakeOKService struct {
	scores []score.View
}

func (f *fakeOKService) SubmitAbsolute(ctx context.Context, user uint, total int) error {
	return nil
}

func (f *fakeOKService) SubmitRelative(ctx context.Context, user uint, variation int) error {
	return nil
}

func (f *fakeOKService) Submit(ctx context.Context, user uint, total int, score int) error {
	return nil
}

func (f *fakeOKService) Find(ctx context.Context, filter interface{}) ([]score.View, error) {
	return f.scores, nil
}

type fakeKOService struct {
}

func (f *fakeKOService) SubmitAbsolute(ctx context.Context, user uint, total int) error {
	return errors.New("error")
}

func (f *fakeKOService) SubmitRelative(ctx context.Context, user uint, variation int) error {
	return errors.New("error")
}

func (f *fakeKOService) Submit(ctx context.Context, user uint, total int, score int) error {
	return nil
}

func (f *fakeKOService) Find(ctx context.Context, filter interface{}) ([]score.View, error) {
	return nil, errors.New("error")
}
