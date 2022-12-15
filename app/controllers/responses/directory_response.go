package responses

import (
	"outclass-api/app/models"
	"outclass-api/app/utils"
	"time"
)

type PostResponse struct {
	Id           string                   `json:"id"`
	ParentId     *string                  `json:"parent_id"`
	OwnerId      string                   `json:"owner_id"`
	ClassroomId  *string                  `json:"classroom_id"`
	Name         string                   `json:"name"`
	Description  *string                  `json:"description"`
	Files        []FileResponse           `json:"files"`
	SharedWith   []UserWithAccessResponse `json:"shared_with"`
	LastModified time.Time                `json:"last_modified"`
	DateCreated  time.Time                `json:"date_created"`
}

type FolderResponse struct {
	Id           string                   `json:"id"`
	ParentId     *string                  `json:"parent_id"`
	OwnerId      string                   `json:"owner_id"`
	ClassroomId  *string                  `json:"classroom_id"`
	Name         string                   `json:"name"`
	Color        *string                  `json:"color"`
	Description  *string                  `json:"description"`
	SharedWith   []UserWithAccessResponse `json:"shared_with"`
	LastModified time.Time                `json:"last_modified"`
	DateCreated  time.Time                `json:"date_created"`
}

type FileResponse struct {
	Name string  `json:"name"`
	Link string  `json:"link"`
	Type *string `json:"type"`
	Size *int64  `json:"size"`
}

type UserWithAccessResponse struct {
	UserId string `json:"user_id"`
	Access string `json:"access"`
}

type FileUploadResponse struct {
	Link string  `json:"link"`
	Type *string `json:"type"`
	Size *int64  `json:"size"`
}

func ToFileUploadResponse(file models.File) FileUploadResponse {
	return FileUploadResponse{
		Link: file.Link,
		Type: file.Type,
		Size: file.Size,
	}
}

func ToPostResponse(directory models.Directory) PostResponse {
	return PostResponse{
		Id:           directory.Id.Hex(),
		ParentId:     utils.ToObjectIdString(*directory.ParentId),
		OwnerId:      directory.OwnerId.Hex(),
		ClassroomId:  utils.ToObjectIdString(*directory.ClassroomId),
		Name:         directory.Name,
		Description:  directory.Description,
		Files:        toFileResponses(directory.Files),
		SharedWith:   toUserWithAccessResponses(directory.SharedWith),
		LastModified: directory.LastModified.Time(),
		DateCreated:  directory.DateCreated.Time(),
	}
}

func ToFolderResponse(directory models.Directory) FolderResponse {
	return FolderResponse{
		Id:           directory.Id.Hex(),
		ParentId:     utils.ToObjectIdString(*directory.ParentId),
		OwnerId:      directory.OwnerId.Hex(),
		ClassroomId:  utils.ToObjectIdString(*directory.ClassroomId),
		Name:         directory.Name,
		Description:  directory.Description,
		Color:        directory.Color,
		SharedWith:   toUserWithAccessResponses(directory.SharedWith),
		LastModified: directory.LastModified.Time(),
		DateCreated:  directory.DateCreated.Time(),
	}
}

func toFileResponses(files []models.File) []FileResponse {
	if len(files) == 0 {
		return nil
	}
	fileResponses := make([]FileResponse, len(files))
	for i, file := range files {
		fileResponses[i].Name = file.Name
		fileResponses[i].Link = file.Link
		fileResponses[i].Type = file.Type
		fileResponses[i].Size = file.Size
	}
	return fileResponses
}

func toUserWithAccessResponses(userWithAccess []models.UserWithAccess) []UserWithAccessResponse {
	if len(userWithAccess) == 0 {
		return nil
	}
	userWithAccessResponses := make([]UserWithAccessResponse, len(userWithAccess))
	for i, userWithAccess := range userWithAccess {
		userWithAccessResponses[i].UserId = userWithAccess.UserId.Hex()
		userWithAccessResponses[i].Access = userWithAccess.Access
	}
	return userWithAccessResponses
}
