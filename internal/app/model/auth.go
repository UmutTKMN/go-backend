package model

import "time"

// Authentication ile ilgili model tanımları

// RegisterRequest kullanıcı kayıt isteklerini temsil eder
type RegisterRequest struct {
	Username    string `json:"username" binding:"required,min=3,max=32"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	BirthDate   string `json:"birth_date"`
	Country     string `json:"country"`
}

// LoginRequest kullanıcı giriş isteklerini temsil eder
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse kimlik doğrulama yanıtını temsil eder
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// TokenClaims token içindeki bilgileri temsil eder
type TokenClaims struct {
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	ExpiredAt time.Time `json:"expired_at"`
}
