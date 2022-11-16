package _auth

import (
	"outclass-api/app/models"
	"outclass-api/app/utils"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserResponse struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func ToTokenResponse(tokens *utils.Tokens) TokenResponse {
	return TokenResponse{
		AccessToken:  tokens.Access,
		RefreshToken: tokens.Refresh,
	}
}

func ToUserResponse(user *models.User) UserResponse {
	return UserResponse{
		Id:    user.Id.Hex(),
		Name:  user.Name,
		Email: user.Email,
	}
}
