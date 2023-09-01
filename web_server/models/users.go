package models

import "time"

type User struct {
	Email            string
	regitration_date time.Time
	password         *string //Only To Be Used with DynamoDB or Postgres
}

type UserModel interface {
	CreateUser(User) User
	deleteUser(string) User
	editUser(User) User
	GetUserByID(string) *User
	ValidatePassword(string, string) User
	GetUserByIDWithCreate(string) (User, bool)
}
