package errorresponse

import (
	"encoding/json"
	"errors"
	"meeting-analyzer/server/commons/utils"
	"net/http"
	"time"

	"meeting-analyzer/server/api/rest/generated"

	"eos2git.cec.lab.emc.com/ISG-Edge/hzp-powerapi-lib-go/powerapi"
)

// errors
var (
	ErrBadPaginationParams       = errors.New("offset must be non-negative and limit must be positive and less than or equal to 1000")
	ErrCheckDriftConflict        = errors.New("another check drift was invoked")
	ErrDeploymentIDNotFound      = errors.New("deployment id not found")
	ErrExecutionIDNotFound       = errors.New("execution id not found")
	ErrInvalidFilterCategory     = errors.New("invalid category")
	ErrInvalidFilterOperator     = errors.New("invalid operator")
	ErrBlueprintRevisionNotFound = errors.New("blueprint revision not found")
)

type ServiceErrorResponse generated.ErrorResponse

func (e *ServiceErrorResponse) Error() string {
	if e.Messages == nil || len(*e.Messages) == 0 {
		return "Invalid error response"
	}
	messages := *e.Messages
	return *messages[0].Message
}

var StrictHTTPServerOptions = generated.StrictHTTPServerOptions{
	RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
		http.Error(w, err.Error(), http.StatusBadRequest)
	},
	ResponseErrorHandlerFunc: ResponseErrorHandlerFunc,
}

func ResponseErrorHandlerFunc(w http.ResponseWriter, _ *http.Request, err error) {
	errorResponse := CreateErrorResponse(err)

	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Set the status code in the response header
	w.WriteHeader(GetStatusCode(errorResponse))

	// Encode the error response and handle potential errors
	if encodeErr := json.NewEncoder(w).Encode(errorResponse); encodeErr != nil {
		http.Error(w, encodeErr.Error(), http.StatusInternalServerError)
	}
}

// CreateErrorResponse creates an ErrorResponse based on the provided error.
// NOTE: If you add a new API, ensure that all possible status codes returned by the backend
// are defined in the openapi specification and handled in this function.
func CreateErrorResponse(err error) *ServiceErrorResponse {
	var errMsg string
	var statusCode generated.HTTPStatusEnum
	var errorResponse *ServiceErrorResponse
	var parseError *powerapi.ParseFilterError

	switch {
	case errors.As(err, &errorResponse):
		return errorResponse
	case errors.Is(err, ErrBadPaginationParams), errors.As(err, &parseError), errors.Is(err, ErrInvalidFilterCategory), errors.Is(err, ErrInvalidFilterOperator):
		errMsg = err.Error()
		statusCode = generated.N400
	case errors.Is(err, ErrDeploymentIDNotFound), errors.Is(err, ErrExecutionIDNotFound), errors.Is(err, ErrBlueprintRevisionNotFound):
		errMsg = "No data found"
		statusCode = generated.N404
	case errors.Is(err, ErrCheckDriftConflict):
		errMsg = "Check drift was already invoked and is in progress"
		statusCode = generated.N409
	default:
		errMsg = err.Error()
		statusCode = generated.N500
	}
	return &ServiceErrorResponse{
		HttpStatusCode: utils.ToPointer(statusCode),
		Messages:       CreateErrorMessages(errMsg),
	}
}

func CreateErrorMessages(errorMessage string) *[]generated.ErrorMessage {
	return &[]generated.ErrorMessage{
		{
			Arguments:   nil,
			Code:        nil,
			Message:     &errorMessage,
			MessageL10n: nil,
			Severity:    utils.ToPointer(generated.ERROR),
			Timestamp:   utils.ToPointer(time.Now()),
		}}
}

func GetStatusCode(errorResponse *ServiceErrorResponse) int {
	if errorResponse != nil && errorResponse.HttpStatusCode != nil {
		return int(*errorResponse.HttpStatusCode)
	}
	return int(generated.N500)
}
