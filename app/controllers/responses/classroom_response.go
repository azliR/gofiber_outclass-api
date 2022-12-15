package responses

import "outclass-api/app/models"

type ClassroomResponse struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	ClassCode   string  `json:"class_code"`
	Description *string `json:"description"`
}

func ToClassroomResponse(classroom models.Classroom) ClassroomResponse {
	return ClassroomResponse{
		Id:          classroom.Id.Hex(),
		Name:        classroom.Name,
		ClassCode:   classroom.ClassCode,
		Description: classroom.Description,
	}
}

func ToClassroomResponses(classrooms []models.Classroom) []ClassroomResponse {
	var classroomResponses []ClassroomResponse
	for _, classroom := range classrooms {
		classroomResponses = append(classroomResponses, ToClassroomResponse(classroom))
	}
	return classroomResponses
}
