package dtos

import (
	"outclass-api/app/models"
)

type CreatePostDto struct {
	ParentId    string  `form:"parent_id" validate:"omitempty,len=24"`
	ClassroomId string  `form:"classroom_id" validate:"omitempty,len=24"`
	Name        string  `form:"name" validate:"required"`
	Description *string `form:"description"`
}

type CreateFolderDto struct {
	ParentId    string  `json:"parent_id" validate:"omitempty,len=24"`
	ClassroomId string  `json:"classroom_id" validate:"omitempty,len=24"`
	Name        string  `json:"name" validate:"required"`
	Color       string  `json:"color" validate:"omitempty,oneof=maraschino cayenne maroon plum eggplant grape orchid lavender carnation strawberry bubblegum magenta salmon tangerine cantaloupe banana lemon honeydew lime spring clover fern moss flora seafoam spindrift teal sky turquoise"`
	Description *string `json:"description"`
}

type UpdatePostDto struct {
	ParentId    string    `json:"parent_id" validate:"omitempty,len=24"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Files       []FileDto `json:"files" validate:"dive,unique=Link"`
}

type FileDto struct {
	Name string  `json:"name" validate:"required"`
	Link string  `json:"link" validate:"required,url"`
	Type *string `json:"type"`
	Size *int64  `json:"size" validate:"omitempty,number"`
}

type UpdateFolderDto struct {
	ParentId    string  `json:"parent_id" validate:"omitempty,len=24"`
	Name        string  `json:"name"`
	Color       *string `json:"color" validate:"omitempty,oneof=null maraschino cayenne maroon plum eggplant grape orchid lavender carnation strawberry bubblegum magenta salmon tangerine cantaloupe banana lemon honeydew lime spring clover fern moss flora seafoam spindrift teal sky turquoise"`
	Description *string `json:"description"`
}

type UserWithAccess struct {
	UserId string `json:"user_id" validate:"required,len=24"`
	Access string `json:"access" validate:"required,oneof=read edit"`
}

type GetDirectoriesDto struct {
	Type        string `query:"type" validate:"required,oneof=folder post"`
	ShareType   string `query:"share_type" validate:"required,oneof=class group personal"`
	ClassroomId string `query:"classroom_id" validate:"required_if=share_type class,omitempty,len=24"`
	ParentId    string `query:"parent_id"`
	Page        uint   `query:"page" validate:"required"`
	PageLimit   uint   `query:"page_limit" validate:"required"`
}

func ToModelFiles(files []FileDto) []models.File {
	modelFiles := make([]models.File, len(files))
	for i, file := range files {
		modelFiles[i].Name = file.Name
		modelFiles[i].Link = file.Link
		modelFiles[i].Type = file.Type
		modelFiles[i].Size = file.Size
	}
	return modelFiles
}
