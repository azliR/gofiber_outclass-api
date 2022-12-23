package responses

import (
	"outclass-api/app/models"
	"outclass-api/app/utils"
	"time"
)

type TaskResponse struct {
	Id           string         `json:"id"`
	CreatorId    string         `json:"owner_id"`
	ClassroomId  *string        `json:"classroom_id"`
	Title        string         `json:"title"`
	Details      *string        `json:"details"`
	Date         time.Time      `json:"date"`
	Repeat       string         `json:"repeat"`
	Color        *string        `json:"color"`
	Files        []FileResponse `json:"files"`
	LastModified time.Time      `json:"last_modified"`
	DateCreated  time.Time      `json:"date_created"`
}

func ToTaskResponse(task models.Task) TaskResponse {
	return TaskResponse{
		Id:           task.Id.Hex(),
		CreatorId:    task.CreatorId.Hex(),
		ClassroomId:  utils.ToObjectIdString(*task.ClassroomId),
		Title:        task.Title,
		Details:      task.Details,
		Date:         task.Date.Time(),
		Repeat:       task.Repeat,
		Color:        &task.Color,
		Files:        ToFileResponses(task.Files),
		LastModified: task.LastModified.Time(),
		DateCreated:  task.DateCreated.Time(),
	}
}

func ToTaskResponses(tasks []models.Task) []TaskResponse {
	taskResponses := make([]TaskResponse, len(tasks))
	for i, task := range tasks {
		taskResponses[i] = ToTaskResponse(task)
	}
	return taskResponses
}
