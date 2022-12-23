package responses

import (
	"outclass-api/app/models"
	"outclass-api/app/utils"
	"time"
)

type EventResponse struct {
	Id           string     `json:"id"`
	CreatorId    string     `json:"owner_id"`
	ClassroomId  *string    `json:"classroom_id"`
	Title        string     `json:"title"`
	StartDate    time.Time  `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	Repeat       string     `json:"repeat"`
	Color        *string    `json:"color"`
	Description  *string    `json:"description"`
	LastModified time.Time  `json:"last_modified"`
	DateCreated  time.Time  `json:"date_created"`
}

func ToEventResponse(event models.Event) EventResponse {
	var endDateTime *time.Time
	if event.EndDate != nil {
		endDate := event.EndDate.Time()
		endDateTime = &endDate
	}
	return EventResponse{
		Id:           event.Id.Hex(),
		CreatorId:    event.CreatorId.Hex(),
		ClassroomId:  utils.ToObjectIdString(*event.ClassroomId),
		Title:        event.Title,
		StartDate:    event.StartDate.Time(),
		EndDate:      endDateTime,
		Repeat:       event.Repeat,
		Color:        &event.Color,
		Description:  event.Description,
		LastModified: event.LastModified.Time(),
		DateCreated:  event.DateCreated.Time(),
	}
}

func ToEventResponses(events []models.Event) []EventResponse {
	eventResponses := make([]EventResponse, len(events))
	for i, event := range events {
		eventResponses[i] = ToEventResponse(event)
	}
	return eventResponses
}
