package model

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type User struct {
	ID                int64  `gorm:"primaryKey" db:"id"`
	UserID            string `gorm:"type:varchar(255)" db:"user_id"`
	Avatar            string `gorm:"type:varchar(255)" db:"avatar"`
	UniqueNickname    string `gorm:"type:varchar" db:"unique_nickname"`
	Email             string `gorm:"type:varchar" db:"email"`
	Level             int    `gorm:"type:int8" db:"level"`
	DeletedTime       int64  `gorm:"type:int8;default:0" db:"deleted_time"`
	CreatedTime       int64  `gorm:"type:int8" db:"created_time"`
	UpdatedTime       int64  `gorm:"type:int8" db:"updated_time"`
	RegisterIP        string `gorm:"type:varchar" db:"register_ip"`
	RegisterUserAgent string `gorm:"type:varchar" db:"register_user_agent"`
	ProfileBackground string `gorm:"type:profile_background" db:"profile_background"`
}

func (u User) SelectOneByUserID(db *gorm.DB, userID string) (*User, error) {
	var entity User
	err := db.Where("user_id = ? AND deleted_time = 0", userID).
		Limit(1).Find(&entity).Error
	if err != nil {
		return nil, err
	}
	if entity.UserID == "" {
		return nil, nil
	}
	return &entity, nil
}

func (u User) SelectOneByEmail(db *gorm.DB, email string) (*User, error) {
	var entity User
	err := db.Where("email = ? AND deleted_time = 0", email).
		Limit(1).Find(&entity).Error
	if err != nil {
		return nil, err
	}
	if entity.UserID == "" {
		return nil, nil
	}
	return &entity, nil
}

type UserPlatformInfo struct {
	UserID string `gorm:"column:user_id;type:character varying(255)" json:"user_id"`
	Email  string `gorm:"column:twitter_id;type:character varying(255)" json:"email"`
}

func (u User) SelectUserWithPlatformByUserId(db *gorm.DB, userId string) (user *UserPlatformInfo, err error) {
	err = db.Raw(`
		SELECT u.user_id, u.email
        FROM users u WHERE u.user_id = ? AND u.deleted_time = 0
	`, userId).Find(&user).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return
}
