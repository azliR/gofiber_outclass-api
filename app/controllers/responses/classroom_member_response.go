package responses

import "outclass-api/app/models"

type ClassroomMemberResponse struct {
	Id            string `json:"id"`
	UserId        string `json:"user_id"`
	ClassroomId   string `json:"classroom_id"`
	StudentId     string `json:"student_id"`
	ClassroomName string `json:"classroom_name"`
	Name          string `json:"name"`
	Role          uint8  `json:"role"`
}

func ToClassroomMemberResponse(classroomMember models.ClassroomMember) ClassroomMemberResponse {
	return ClassroomMemberResponse{
		Id:            classroomMember.Id.Hex(),
		UserId:        classroomMember.UserId.Hex(),
		ClassroomId:   classroomMember.ClassroomId.Hex(),
		StudentId:     classroomMember.StudentId,
		ClassroomName: classroomMember.ClassroomName,
		Name:          classroomMember.Name,
		Role:          classroomMember.Role,
	}
}

func ToClassroomMemberResponses(classroomMembers []models.ClassroomMember) []ClassroomMemberResponse {
	var classroomMemberResponses []ClassroomMemberResponse
	for _, classroomMember := range classroomMembers {
		classroomMemberResponses = append(classroomMemberResponses, ToClassroomMemberResponse(classroomMember))
	}
	return classroomMemberResponses
}
