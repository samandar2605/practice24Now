package models

type User struct {
	Id          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Email       string `json:"email" db:"email"`
	Password    string `json:"password"`
	AccessToken string `json:"access_token"`
}

type CreateUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type VerifyUser struct {
	Email            string `json:"email"`
	VerificationCode int    `json:"verification_code"`
}

type RedisUser struct {
	Name             string `json:"name"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	VerificationCode int    `json:"verification_code"`
}
