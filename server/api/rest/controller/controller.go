// Copyright © 2022 Dell Inc. or its subsidiaries. All Rights Reserved.
package controller

import (
	"context"
	"encoding/json"
	"meeting-analyzer/server/api/rest/generated"
	"meeting-analyzer/server/models"
	"meeting-analyzer/server/services/service"
)

type controller struct {
	svc service.Service
}

func NewController(svc service.Service) generated.StrictServerInterface {
	return &controller{
		svc: svc,
	}
}

func (c *controller) GetMeetingSummaries(ctx context.Context, request generated.GetMeetingSummariesRequestObject) (generated.GetMeetingSummariesResponseObject, error) {
	panic("implement me")
}

func (c *controller) GenerateMeetingSummary(ctx context.Context, request generated.GenerateMeetingSummaryRequestObject) (generated.GenerateMeetingSummaryResponseObject, error) {
	var meetingDetails models.MeetingDetails
	bytes, _ := json.Marshal(request.Body)
	json.Unmarshal(bytes, &meetingDetails)
	res, _ := c.svc.GenerateMeetingSummary(ctx, &meetingDetails)
	return generated.GenerateMeetingSummary202JSONResponse(*res), nil
}
func (c *controller) GetMeetingSummaryById(ctx context.Context, request generated.GetMeetingSummaryByIdRequestObject) (generated.GetMeetingSummaryByIdResponseObject, error) {
	panic("implement me")
}
