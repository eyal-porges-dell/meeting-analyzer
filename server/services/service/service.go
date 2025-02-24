// Copyright © 2022 Dell Inc. or its subsidiaries. All Rights Reserved.

package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"meeting-analyzer/server/api/rest/generated"
	"meeting-analyzer/server/models"
	"meeting-analyzer/server/repositories"
	"net/http"

	log "eos2git.cec.lab.emc.com/ISG-Edge/hzp-go-commons/logger"
)

type Service interface {
	GenerateMeetingSummary(ctx context.Context, meetingDetails *models.MeetingDetails) (*generated.GenerateMeetingSummaryResponse, error)
}

type svc struct {
}

func NewSvc(ctx context.Context, repo repositories.Repository) (Service, error) { return &svc{}, nil }

func (s *svc) GenerateMeetingSummary(ctx context.Context, meetingDetails *models.MeetingDetails) (*generated.GenerateMeetingSummaryResponse, error) {
	s.callAI(ctx, meetingDetails)
	return &generated.GenerateMeetingSummaryResponse{}, nil
}

func (s *svc) callAI(ctx context.Context, meetingDetails *models.MeetingDetails) {
	// Define the URL of the API endpoint
	url := "https://chat.dell.com/api/chat/completions"

	// Define the bearer token for authorization
	bearerToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImE0ZjIzYjYxLWZmNDUtNGQ2NC05MGM1LTdiOGIxNmZkOTRkMCJ9.NYm8gN3unrsT5o20w6uIz9VTJMJvmQrwAgkmLcJo7MY"

	// Define the request payload
	payload := map[string]interface{}{
		"stream": false,
		"model":  "4df3c44b-b22d-4c14-bed3-2a7542f74765",
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": formatTranscription(meetingDetails),
			},
		},
		"session_id": "lMPhtvj7-bhj04RiAIlD",
		"chat_id":    "a78bd7fa-a668-4562-b833-d4073381691a",
		"id":         "b9f51806-abf8-424e-ae58-3e13c147166a",
	}

	// Marshal the payload into a JSON byte slice
	requestBody, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(ctx, nil, "", err, "Failed to marshal the request payload")
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatal(ctx, nil, "", err, "Failed to create the request")
	}

	// Set the Authorization header
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	// Create an HTTP client and make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(ctx, nil, "", err, "Failed to make the GET request")
	}
	// Make sure to close the response body at the end
	defer resp.Body.Close()

	// Check for a successful response status code
	if resp.StatusCode != http.StatusOK {
		log.Fatal(ctx, nil, "", err, "Unexpected status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(ctx, nil, "", err, "Failed to read the response body")
	}

	// Define a structure to match the expected JSON response
	var result map[string]interface{}

	// Unmarshal the JSON response into the structure
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatal(ctx, nil, "", err, "Failed to unmarshal JSON response")
	}

	// Print the JSON response
	fmt.Println("Response JSON:", result)
}

// Helper function to format the transcription content
func formatTranscription(meetingDetails *models.MeetingDetails) string {
	var content string
	for _, t := range meetingDetails.Transcription {
		content += fmt.Sprintf("\"%s\",\"%s\"\n\"%s\"\n", t.Member, t.Timestamp, t.Content)
	}
	return fmt.Sprintf("Meeting Transcription: %s\n%s", meetingDetails.MeetingTitle, content)
}
