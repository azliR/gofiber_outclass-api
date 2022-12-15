package dtos

import (
	"errors"
	"outclass-api/app/constants"
)

// type CreateClassroomMemberDto struct {
// 	StudentId string `json:"student_id" validate:"required"`
// }

type GetClassroomMembersDto struct {
	Page      uint `query:"page" validate:"required"`
	PageLimit uint `query:"page_limit" validate:"required"`
}

func ToModelRole(role string) (uint8, error) {
	switch role {
	case constants.OwnerClassroomRoleName:
		return 1, nil
	case constants.AdminClassroomRoleName:
		return 2, nil
	case constants.MemberClassroomRoleName:
		return 3, nil
	default:
		return 0, errors.New("role is not valid")
	}
}

func FromModelRole(role uint8) (string, error) {
	switch role {
	case 1:
		return constants.OwnerClassroomRoleName, nil
	case 2:
		return constants.AdminClassroomRoleName, nil
	case 3:
		return constants.MemberClassroomRoleName, nil
	default:
		return "", errors.New("role is not valid")
	}
}
