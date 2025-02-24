package errorresponse

import (
	"encoding/json"
	"errors"
	"meeting-analyzer/server/api/rest/generated"
	"meeting-analyzer/server/commons/utils"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestResponseErrorHandlerFunc(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectedBody   *generated.ErrorResponse
	}{
		{
			name:           "BadPaginationParams",
			err:            ErrBadPaginationParams,
			expectedStatus: http.StatusBadRequest,
			expectedBody: &generated.ErrorResponse{
				HttpStatusCode: utils.ToPointer(generated.N400),
				Messages: &[]generated.ErrorMessage{
					{
						Message:   utils.ToPointer("offset must be non-negative and limit must be positive and less than or equal to 1000"),
						Severity:  utils.ToPointer(generated.ERROR),
						Timestamp: utils.ToPointer(time.Now()),
					},
				},
			},
		},
		{
			name:           "DeploymentNotFound",
			err:            ErrDeploymentIDNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody: &generated.ErrorResponse{
				HttpStatusCode: utils.ToPointer(generated.N404),
				Messages: &[]generated.ErrorMessage{
					{
						Message:   utils.ToPointer("No data found"),
						Severity:  utils.ToPointer(generated.ERROR),
						Timestamp: utils.ToPointer(time.Now()),
					},
				},
			},
		},
		{
			name:           "CheckDriftConflict",
			err:            ErrCheckDriftConflict,
			expectedStatus: http.StatusConflict,
			expectedBody: &generated.ErrorResponse{
				HttpStatusCode: utils.ToPointer(generated.N409),
				Messages: &[]generated.ErrorMessage{
					{
						Message:   utils.ToPointer("Check drift was already invoked and is in progress"),
						Severity:  utils.ToPointer(generated.ERROR),
						Timestamp: utils.ToPointer(time.Now()),
					},
				},
			},
		},
		{
			name:           "InternalServerError",
			err:            errors.New("internal server error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody: &generated.ErrorResponse{
				HttpStatusCode: utils.ToPointer(generated.N500),
				Messages: &[]generated.ErrorMessage{
					{
						Message: utils.ToPointer("internal server error"),
					},
				},
			},
		},
		{
			name:           "InvalidFilterCategory",
			err:            ErrInvalidFilterCategory,
			expectedStatus: http.StatusBadRequest,
			expectedBody: &generated.ErrorResponse{
				HttpStatusCode: utils.ToPointer(generated.N400),
				Messages: &[]generated.ErrorMessage{
					{
						Message:   utils.ToPointer("invalid category"),
						Severity:  utils.ToPointer(generated.ERROR),
						Timestamp: utils.ToPointer(time.Now()),
					},
				},
			},
		},
		{
			name:           "InvalidFilterOperator",
			err:            ErrInvalidFilterOperator,
			expectedStatus: http.StatusBadRequest,
			expectedBody: &generated.ErrorResponse{
				HttpStatusCode: utils.ToPointer(generated.N400),
				Messages: &[]generated.ErrorMessage{
					{
						Message:   utils.ToPointer("invalid operator"),
						Severity:  utils.ToPointer(generated.ERROR),
						Timestamp: utils.ToPointer(time.Now()),
					},
				},
			},
		},
		{
			name:           "BlueprintRevisionNotFound",  // New test case
			err:            ErrBlueprintRevisionNotFound, // New test case
			expectedStatus: http.StatusNotFound,          // New test case
			expectedBody: &generated.ErrorResponse{
				HttpStatusCode: utils.ToPointer(generated.N404), // New test case
				Messages: &[]generated.ErrorMessage{
					{
						Message:   utils.ToPointer("No data found"), // New test case
						Severity:  utils.ToPointer(generated.ERROR), // New test case
						Timestamp: utils.ToPointer(time.Now()),      // New test case
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			ResponseErrorHandlerFunc(recorder, nil, tt.err)

			assert.Equal(t, tt.expectedStatus, recorder.Code)

			var actualBody generated.ErrorResponse
			err := json.Unmarshal(recorder.Body.Bytes(), &actualBody)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedBody.HttpStatusCode, actualBody.HttpStatusCode)
			assert.Equal(t, (*tt.expectedBody.Messages)[0].Message, (*actualBody.Messages)[0].Message)
		})
	}
}

func TestCreateErrorMessages(t *testing.T) {
	errorMessage := "Sample error message"
	expectedMessages := &[]generated.ErrorMessage{
		{
			Message:   &errorMessage,
			Severity:  utils.ToPointer(generated.ERROR),
			Timestamp: utils.ToPointer(time.Now()),
		},
	}

	actualMessages := CreateErrorMessages(errorMessage)
	assert.Equal(t, (*expectedMessages)[0].Message, (*actualMessages)[0].Message)
	assert.Equal(t, generated.ERROR, *(*actualMessages)[0].Severity)
	assert.WithinDuration(t, time.Now(), *(*actualMessages)[0].Timestamp, time.Second)
}

func TestGetStatusCode_NilPointer(t *testing.T) {
	actualStatusCode := GetStatusCode(nil)
	assert.Equal(t, http.StatusInternalServerError, actualStatusCode)
}
