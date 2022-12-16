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

type ClassroomMemberProfileResponse struct {
	Id                   string  `json:"id"`
	UserId               string  `json:"user_id"`
	ClassroomId          string  `json:"classroom_id"`
	StudentId            string  `json:"student_id"`
	Name                 string  `json:"name"`
	Role                 uint8   `json:"role"`
	ClassroomName        string  `json:"classroom_name"`
	ClassCode            string  `json:"class_code"`
	ClassroomDescription *string `json:"description"`
	ClassMembersCount    int64   `json:"class_members_count"`
}

func ToClassroomMemberResponse(classroomMember models.ClassroomMember) ClassroomMemberResponse {
	return ClassroomMemberResponse{
		Id:            classroomMember.Id.Hex(),
		UserId:        classroomMember.UserId.Hex(),
		ClassroomId:   classroomMember.ClassroomId.Hex(),
		StudentId:     classroomMember.StudentId,
		Name:          classroomMember.Name,
		Role:          classroomMember.Role,
		ClassroomName: classroomMember.ClassroomName,
	}
}

func ToClassroomMemberProfileResponse(classroom models.Classroom, classroomMember models.ClassroomMember, classMembersCount int64) ClassroomMemberProfileResponse {
	return ClassroomMemberProfileResponse{
		Id:                   classroomMember.Id.Hex(),
		UserId:               classroomMember.UserId.Hex(),
		ClassroomId:          classroomMember.ClassroomId.Hex(),
		StudentId:            classroomMember.StudentId,
		Name:                 classroomMember.Name,
		Role:                 classroomMember.Role,
		ClassroomName:        classroomMember.ClassroomName,
		ClassCode:            classroom.ClassCode,
		ClassroomDescription: classroom.Description,
		ClassMembersCount:    classMembersCount,
	}
}

func ToClassroomMemberResponses(classroomMembers []models.ClassroomMember) []ClassroomMemberResponse {
	var classroomMemberResponses []ClassroomMemberResponse
	for _, classroomMember := range classroomMembers {
		classroomMemberResponses = append(classroomMemberResponses, ToClassroomMemberResponse(classroomMember))
	}
	return classroomMemberResponses
}
