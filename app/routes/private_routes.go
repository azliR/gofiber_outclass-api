package routes

import (
	"outclass-api/app/handlers"
	"outclass-api/app/middleware"

	"github.com/gofiber/fiber/v2"
)

func PrivateRoutes(a *fiber.App) {
	route := a.Group("/api/v1")

	route.Get("/user/profile", middleware.JWTProtected(), handlers.UserProfile)
	route.Post("/user/sign/out", middleware.JWTProtected(), handlers.UserSignOut)
	route.Post("/token/renew", middleware.JWTProtected(), handlers.RenewTokens)

	route.Get("/files/get/:fileId", middleware.JWTProtected(), handlers.GetFile)
	route.Post("/files/upload", middleware.JWTProtected(), handlers.UploadFile)
	route.Delete("/files/:fileId", middleware.JWTProtected(), handlers.DeleteFile)

	route.Post("/directories/post", middleware.JWTProtected(), handlers.CreatePost)
	route.Put("/directories/post", middleware.JWTProtected(), handlers.UpdatePostById)

	route.Post("/directories/folder", middleware.JWTProtected(), handlers.CreateFolder)
	route.Put("/directories/folder", middleware.JWTProtected(), handlers.UpdateFolderById)

	route.Get("/directories", middleware.JWTProtected(), handlers.GetDirectoriesByParentId)
	route.Get("/directories/:directoryId", middleware.JWTProtected(), handlers.GetDirectory)
	route.Delete("/directories/:directoryId", middleware.JWTProtected(), handlers.DeleteDirectory)
}
