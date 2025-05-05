package models

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// RegisterResponse represents the response for user registration
type RegisterResponse struct {
	Message string `json:"message"`
}

// TokenResponse represents the response containing access and refresh tokens
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// ProfileResponse represents the response for user profile
type ProfileResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// NewErrorResponse creates a new error response
func NewErrorResponse(err string) ErrorResponse {
	return ErrorResponse{Error: err}
}

// NewRegisterResponse creates a new registration response
func NewRegisterResponse() RegisterResponse {
	return RegisterResponse{Message: "User created successfully"}
}

// NewTokenResponse creates a new token response
func NewTokenResponse(accessToken, refreshToken string, expiresIn int64) TokenResponse {
	return TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
	}
}

// NewProfileResponse creates a new profile response
func NewProfileResponse(username, email string) ProfileResponse {
	return ProfileResponse{
		Username: username,
		Email:    email,
	}
}
