package management

import (
	logger "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"one-day-server/internal/db/mysql"
)

var (
	users = make(map[string]*User)
)

func GetUserByAPIKey(apiKey string) *User {
	if userInMem, ok := users[apiKey]; !ok {
		user := &User{}
		if err := mysql.DB().Where("api_key = ?", apiKey).First(user).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				logger.Errorf("get %s api sercret failed, err: %s, obtain user from CORE", apiKey, err)
			} else {
				logger.Infof("get empty user: %s, obtain user from CORE", apiKey)
			}
			//secret = getUserAPISecretFromCore(apiKey)
			if user.Secret == "" {
				return nil
			}
		}
		users[user.APIKey] = user
		return user
	} else {
		return userInMem
	}
}

func GetUserAddressByAPIKey(apiKey string) string {
	return users[apiKey].Address
}

type UserSignData struct {
	Secret    string
	Body      string
	Timestamp int64
}

type User struct {
	Id        int64         `gorm:"column:user_id"`    //
	Address   string        `gorm:"column:address"`    //
	VesselKey string        `gorm:"column:vessel_key"` //
	APIKey    string        `gorm:"column:api_key"`    //
	Secret    string        `gorm:"column:secret"`     //
	SignData  *UserSignData `gorm:"-:all"`             // only be used in memory
}

func (m *User) TableName() string {
	return "users"
}

func AddUser(user *User) error {

	return nil
}
