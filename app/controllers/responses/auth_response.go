package responses

import (
	"outclass-api/app/models"
	"outclass-api/app/utils"
)

type CreateUserResponse struct {
	User  UserResponse  `json:"user"`
	Token TokenResponse `json:"token"`
}

type UserResponse struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type TokenResponse struct {
	AccessToken           string `json:"access_token"`
	TokenExpiresIn        string `json:"token_expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn string `json:"refresh_token_expires_in"`
}

func ToCreateUserResponse(user models.User, tokens utils.Tokens) CreateUserResponse {
	return CreateUserResponse{
		User:  ToUserResponse(user),
		Token: ToTokenResponse(tokens),
	}
}

func ToUserResponse(user models.User) UserResponse {
	return UserResponse{
		Id:    user.Id.Hex(),
		Name:  user.Name,
		Email: user.Email,
	}
}

func ToTokenResponse(tokens utils.Tokens) TokenResponse {
	return TokenResponse{
		AccessToken:           tokens.Access,
		TokenExpiresIn:        tokens.TokenExpiresIn,
		RefreshToken:          tokens.Refresh,
		RefreshTokenExpiresIn: tokens.RefreshTokenExpiresIn,
	}
}
