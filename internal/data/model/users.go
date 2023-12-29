package model

import (
	"gorm.io/gorm"
)

type User struct {
	ID                int64   `gorm:"primaryKey"`
	AppID             string  `gorm:"type:varchar(255)"`
	UserID            string  `gorm:"type:varchar(255)"`
	Avatar            string  `gorm:"type:varchar(255)"`
	UniqueNickname    string  `gorm:"type:varchar"`
	Email             *string `gorm:"type:varchar"`
	Level             int     `gorm:"type:int8"`
	Wxp               int     `gorm:"type:int8"`
	NextLevelWxp      int     `gorm:"type:int8"`
	WowBalance        float64 `gorm:"type:numeric"`
	Dragonball        int     `gorm:"type:int8"`
	SuspectBot        bool    `gorm:"type:bool"`
	DeletedTime       int64   `gorm:"type:int8;default:0"`
	CreatedTime       int64   `gorm:"type:int8"`
	UpdatedTime       int64   `gorm:"type:int8"`
	ShowGuideProfile  bool    `gorm:"type:bool"`
	ShowGuideGames    bool    `gorm:"type:bool"`
	ShowGuideEvents   bool    `gorm:"type:bool"`
	RegisterIP        string  `gorm:"type:varchar"`
	RegisterUserAgent string  `gorm:"type:varchar"`
	ProfileBackground string  `gorm:"type:profile_background"`
	AvatarType        int8    `gorm:"type:avatar_type"`
	AvatarFrame       string  `gorm:"type:avatar_frame"`
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

type UserIdModel struct {
	UserID string `gorm:"column:user_id;type:character varying(255)" json:"user_id"`
}

func (u User) SelectUserIdExistBySql(db *gorm.DB, userId, sql string) (exist bool, err error) {
	var userIdModel *UserIdModel
	err = db.Raw("(?) INTERSECT (select ? as user_id)", sql, userId).Find(&userIdModel).Error
	if err != nil {
		return
	}

	exist = userIdModel.UserID != ""
	return
}

type UserPlatformInfo struct {
	UserID      string   `gorm:"column:user_id;type:character varying(255)" json:"user_id"`
	SteamID     string   `gorm:"column:steam_id;type:character varying(255)" json:"steam_id"`
	DiscordID   string   `gorm:"column:discord_id;type:character varying(255)" json:"discord_id"`
	TwitterID   string   `gorm:"column:twitter_id;type:character varying(255)" json:"twitter_id"`
	XboxID      string   `gorm:"column:xbox_id;type:character varying(255)" json:"xbox_id"`
	EpicID      string   `gorm:"column:epic_id;type:character varying(255)" json:"epic_id"`
	Email       string   `gorm:"column:twitter_id;type:character varying(255)" json:"email"`
	SuspectBot  bool     `gorm:"column:user_id;type:boolean" json:"suspect_bot"`
	WalletAddrs []string `gorm:"column:wallet_addrs;type:jsonb;serializer:json" json:"wallet_addrs"`
}

func (u User) SelectUserWithPlatformByUserId(db *gorm.DB, userId string) (user *UserPlatformInfo, err error) {
	err = db.Raw(`
		SELECT u.suspect_bot,u.user_id,(
            SELECT json_agg(wallet_addr) FROM user_wallets WHERE user_id = u.user_id and deleted_time=0
            ) wallet_addrs,
            (SELECT discord_id FROM user_discords WHERE user_id=u.user_id AND deleted_time=0) discord_id,
            (SELECT twitter_id FROM user_twitters WHERE user_id=u.user_id AND deleted_time=0) twitter_id,
            (SELECT steam_id FROM user_steams WHERE user_id=u.user_id AND deleted_time=0) steam_id,
            (SELECT xbox_id FROM user_xboxes WHERE user_id=u.user_id AND deleted_time=0) xbox_id,
            (SELECT epic_id FROM user_epic WHERE user_id=u.user_id AND deleted_time=0) epic_id,
            u.email
        FROM users u WHERE u.user_id = ? AND u.deleted_time = 0
	`, userId).Find(&user).Error
	if err == nil && user.WalletAddrs == nil {
		user.WalletAddrs = []string{}
	}

	return
}
