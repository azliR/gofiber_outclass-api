package routes

import (
	"outclass-api/app/controllers"
	"outclass-api/app/middleware"

	"github.com/gofiber/fiber/v2"
)

func PrivateRoutes(a *fiber.App) {
	route := a.Group("/api/v1")

	route.Get("/user/profile", middleware.JWTProtected(), controllers.UserProfile)
	route.Post("/user/sign/out", middleware.JWTProtected(), controllers.UserSignOut)
	route.Post("/token/renew", middleware.JWTProtected(), controllers.RenewTokens)

	route.Get("/files/get/:fileId", middleware.JWTProtected(), controllers.GetFile)
	route.Post("/files/upload", middleware.JWTProtected(), controllers.UploadFile)
	route.Delete("/files/:fileId", middleware.JWTProtected(), controllers.DeleteFile)

	route.Post("/directories/post", middleware.JWTProtected(), controllers.CreatePost)
	route.Put("/directories/post", middleware.JWTProtected(), controllers.UpdatePostById)

	route.Post("/directories/folder", middleware.JWTProtected(), controllers.CreateFolder)
	route.Put("/directories/folder", middleware.JWTProtected(), controllers.UpdateFolderById)

	route.Get("/directories", middleware.JWTProtected(), controllers.GetDirectoriesByParentId)
	route.Get("/directories/:directoryId", middleware.JWTProtected(), controllers.GetDirectory)
	route.Delete("/directories/:directoryId", middleware.JWTProtected(), controllers.DeleteDirectory)
}
