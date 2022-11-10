package _directory

// type DirectoryResponse struct {
// 	Id           string                   `json:"id" validate:"required"`
// 	ParentId     string                   `json:"parent_id"`
// 	OwnerId      string                   `json:"owner_id" validate:"required"`
// 	ClassroomId  string                   `json:"classroom_id"`
// 	Name         string                   `json:"name" validate:"required"`
// 	Type         string                   `json:"type" validate:"required,oneof=folder post"`
// 	Description  string                   `json:"description"`
// 	Files        []FileResponse           `json:"files" validate:"required_unless=Type folder,dive"`
// 	SharedWith   []UserWithAccessResponse `json:"shared_with" validate:"dive"`
// 	LastModified string                   `json:"last_modified" validate:"required"`
// 	DateCreated  string                   `json:"date_created" validate:"required"`
// }

// type FileResponse struct {
// 	Link string `json:"link" validate:"required,url"`
// 	Type string `json:"type" validate:"required"`
// 	Size int64  `json:"size" validate:"required,number"`
// }

type UserWithAccessResponse struct {
	UserId string `json:"user_id" validate:"required"`
	Access string `json:"access" validate:"required,oneof=read edit"`
}

type FileUploadResponse struct {
	FileName string `json:"file_name"`
	Link     string `json:"link"`
}

// func ToDirectoryResponse(directory models.Directory) DirectoryResponse {
// 	fileResponses := make([]FileResponse, len(directory.Files))
// 	for i, file := range directory.Files {
// 		fileResponses[i].Link = file.Link
// 		fileResponses[i].Type = file.Type
// 		fileResponses[i].Size = file.Size
// 	}
// 	userWithAccessResponses := make([]UserWithAccessResponse, len(directory.SharedWith))
// 	for i, userWithAccess := range directory.SharedWith {
// 		userWithAccessResponses[i].UserId = userWithAccess.UserId.Hex()
// 		userWithAccessResponses[i].Access = userWithAccess.Access
// 	}

// 	return DirectoryResponse{
// 		Id:           directory.Id.Hex(),
// 		ParentId:     directory.ParentId.Hex(),
// 		OwnerId:      directory.OwnerId.Hex(),
// 		ClassroomId:  directory.ClassroomId.Hex(),
// 		Name:         directory.Name,
// 		Type:         directory.Type,
// 		Description:  directory.Description,
// 		Files:        fileResponses,
// 		SharedWith:   userWithAccessResponses,
// 		LastModified: directory.LastModified.Time().String(),
// 		DateCreated:  directory.DateCreated.Time().String(),
// 	}
// }
