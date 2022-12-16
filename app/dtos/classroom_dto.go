package dtos

type CreateClassroomDto struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description"`
	StudentId   string  `json:"student_id" validate:"required"`
}

type JoinClassroomDto struct {
	ClassCode string `json:"class_code" validate:"required,len=16"`
	StudentId string `json:"student_id" validate:"required"`
}

type UpdateClassroomDto struct {
	Name        string  `json:"name" `
	Description *string `json:"description"`
}
