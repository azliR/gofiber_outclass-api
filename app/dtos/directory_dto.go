package dtos

import "outclass-api/app/models"

type CreatePostDto struct {
	ParentId    string    `json:"parent_id"`
	ClassroomId string    `json:"classroom_id"`
	Name        string    `json:"name" validate:"required"`
	Description *string   `json:"description"`
	Files       []FileDto `json:"files" validate:"required,dive,unique=Link"`
}

type FileDto struct {
	Link string `json:"link" validate:"required,url"`
	Type string `json:"type" validate:"required"`
	Size int64  `json:"size" validate:"required,number"`
}

type CreateFolderDto struct {
	ParentId    string  `json:"parent_id" validate:"omitempty,len=24"`
	ClassroomId string  `json:"classroom_id" validate:"required,omitempty,len=24"`
	Name        string  `json:"name" validate:"required"`
	Color       string  `json:"color" validate:"omitempty,oneof=maraschino cayenne maroon plum eggplant grape orchird lavender carnation strawberry bubblegum magenta salmon tangerine cantaloupe banana lemon honeydew lime spring clover fern moss flora sea foam spindrift teal sky turquoise"`
	Description *string `json:"description"`
}

type UpdatePostDto struct {
	Id          string    `json:"id" validate:"required,len=24"`
	ParentId    string    `json:"parent_id" validate:"omitempty,len=24"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Files       []FileDto `json:"files" validate:"dive,unique=Link"`
}

type UpdateFolderDto struct {
	Id          string `json:"id" validate:"required,len=24"`
	ParentId    string `json:"parent_id" validate:"omitempty,len=24"`
	Name        string `json:"name"`
	Color       string `json:"color" validate:"omitempty,oneof=null maraschino cayenne maroon plum eggplant grape orchird lavender carnation strawberry bubblegum magenta salmon tangerine cantaloupe banana lemon honeydew lime spring clover fern moss flora sea foam spindrift teal sky turquoise"`
	Description string `json:"description"`
}

type UserWithAccess struct {
	UserId string `json:"user_id" validate:"required,len=24"`
	Access string `json:"access" validate:"required,oneof=read edit"`
}

type GetDirectoriesDto struct {
	ParentId  string `query:"parent_id" validate:"required"`
	Page      uint   `query:"page" validate:"required"`
	PageLimit uint   `query:"page_limit" validate:"required"`
}

func ToModelFiles(files []FileDto) []models.File {
	modelFiles := make([]models.File, len(files))
	for i, file := range files {
		modelFiles[i].Link = file.Link
		modelFiles[i].Type = file.Type
		modelFiles[i].Size = file.Size
	}
	return modelFiles
}
