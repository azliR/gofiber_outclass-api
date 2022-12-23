package routes

import (
	"outclass-api/app/controllers"
	"outclass-api/app/middleware"

	"github.com/gofiber/fiber/v2"
)

func PrivateRoutes(a *fiber.App) {
	route := a.Group("/api/v1")

	route.Put("/user", middleware.JWTProtected(), controllers.UpdateUser)
	route.Get("/user/profile", middleware.JWTProtected(), controllers.UserProfile)
	route.Get("/user/classrooms", middleware.JWTProtected(), controllers.GetClassroomMembersByUserId)
	route.Post("/user/sign/out", middleware.JWTProtected(), controllers.UserSignOut)
	route.Post("/user/refresh", middleware.JWTProtected(), controllers.RenewTokens)

	route.Post("/classrooms", middleware.JWTProtected(), controllers.CreateClassroom)
	route.Post("/classrooms/join", middleware.JWTProtected(), controllers.JoinClassroomByClassCode)
	route.Put("/classrooms/:classroomId", middleware.JWTProtected(), controllers.UpdateClassroomById)
	route.Get("/classrooms/:classroomId", middleware.JWTProtected(), controllers.GetClassroomById)
	route.Get("/classrooms/code/:classCode", middleware.JWTProtected(), controllers.GetClassroomByClassCode)
	route.Delete("/classrooms/:classroomId", middleware.JWTProtected(), controllers.DeleteClassroom)

	// route.Post("/classrooms/:classroomId/members", middleware.JWTProtected(), controllers.CreateClassroomMember)
	route.Get("/classrooms/:classroomId/members", middleware.JWTProtected(), controllers.GetClassroomMembersByClassroomId)
	route.Get("/classrooms/:classroomId/members/profile", middleware.JWTProtected(), controllers.GetClassroomMemberByUserIdAndClassroomId)

	// route.Post("/files", middleware.JWTProtected(), controllers.UploadFile)
	route.Get("/files/:fileId", middleware.JWTProtected(), controllers.GetFile)
	// route.Delete("/files/:fileId", middleware.JWTProtected(), controllers.DeleteFile)

	route.Post("/directories/posts", middleware.JWTProtected(), controllers.CreatePost)
	route.Put("/directories/posts/:postId", middleware.JWTProtected(), controllers.UpdatePostById)

	route.Post("/directories/folders", middleware.JWTProtected(), controllers.CreateFolder)
	route.Put("/directories/folders/:folderId", middleware.JWTProtected(), controllers.UpdateFolderById)

	route.Get("/directories", middleware.JWTProtected(), controllers.GetDirectoriesByParentId)
	route.Get("/directories/:directoryId", middleware.JWTProtected(), controllers.GetDirectoryById)
	route.Delete("/directories/:directoryId", middleware.JWTProtected(), controllers.DeleteDirectory)

	route.Post("/events", middleware.JWTProtected(), controllers.CreateEvent)
	route.Get("/events/class/:classroomId", middleware.JWTProtected(), controllers.GetEventsByClassroomId)

	route.Post("/tasks", middleware.JWTProtected(), controllers.CreateTask)
	route.Get("/tasks/class/:classroomId", middleware.JWTProtected(), controllers.GetTasksByClassroomId)
}
