package handler

import (
	"net/http"
	"testing"
)

func TestUpdateCategory(t *testing.T) {
	tests := []struct {
		name           string
		param          string
		inputJSON      string
		mockOn         int
		mockReturn     error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid Request",
			param:          "1",
			inputJSON:      `{"name": "Name", "description": "Description"}`,
			mockOn:         1,
			mockReturn:     nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"Request Status": "Changes completed"}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			
		})
	}
}
