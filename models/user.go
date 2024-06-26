package models

import (
	"errors"
	"final-project/helpers"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type User struct {
	GormModel
	Username        string `json:"username" gorm:"unique;not null" form:"username" valid:"required~Username is required"`
	Email           string `json:"email" gorm:"unique;not null" form:"email" valid:"required~Email is required, email~Email is invalid"`
	Password        string `json:"password" gorm:"not null" form:"password" valid:"required~Password is required, minstringlength(6)~Password must be at least 6 characters"`
	ProfileImageURL string `json:"profile_image_url" form:"profile_image_url" valid:"required~Profile Image URL is required, url~Invalid URL format"`
	Age             int    `json:"age" gorm:"not null" form:"age" valid:"required~Age is required, range(8|100)~Age must be at least 8"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(u)

	if errCreate != nil {
		err = errCreate
		return
	}

	u.Password = helpers.HashPassword(u.Password)

	err = nil
	return
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	if u.Email == "" && u.Username == "" {
		err = errors.New("Email and Username is required")
		return
	} else if u.Username == "" {
		err = errors.New("Username is required")
		return
	} else if u.Email == "" {
		err = errors.New("Email is required")
		return
	}

	err = nil
	return
}
