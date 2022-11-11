package dtos

type CreateClassroomDto struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description"`
}

type UpdateClassroomDto struct {
	Name        string  `json:"name" `
	Description *string `json:"description"`
}
