package models

type Transcription struct {
	Member    string `json:"member"`
	Timestamp string `json:"timestamp"`
	Content   string `json:"content"`
}

type MeetingDetails struct {
	MeetingID     string          `json:"meeting_id"`
	MeetingTitle  string          `json:"meeting_title"`
	Transcription []Transcription `json:"transcription"`
}
