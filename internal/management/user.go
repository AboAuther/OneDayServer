package management

import (
	"errors"
	"fmt"

	"one-day-server/internal/db/mysql"
)

var (
	usernameToUser = make(map[string]*User)
	phoneToUser    = make(map[string]*User)
	uidToUser      = make(map[int64]*User)

	UserNotExist = errors.New("user not exist")
)

func Init() {
	users := []*User{}
	if err := mysql.DB().Debug().Model(&User{}).Find(&users).Error; err != nil {
		panic(err)
	}
	for _, user := range users {
		usernameToUser[user.Username] = user
		phoneToUser[user.Phone] = user
		uidToUser[user.Id] = user
	}
}

type User struct {
	Id           int64  `gorm:"column:user_id"` //
	Username     string `gorm:"column:username"`
	Password     string `gorm:"column:password"`
	Email        string `gorm:"column:email"`
	Phone        string `gorm:"column:phone"`
	Gender       string `gorm:"column:gender"`
	Age          int    `gorm:"column:age"`
	IsVip        bool   `gorm:"column:is_vip"`
	RefreshToken string `gorm:"column:refresh_token"`
}

func (m *User) TableName() string {
	return "users"
}

func AddUser(user *User) error {
	if err := mysql.DB().Model(&User{}).Debug().Create(user).Error; err != nil {
		return fmt.Errorf("create user failed, err: %s", err)
	}
	usernameToUser[user.Username] = user
	phoneToUser[user.Phone] = user
	uidToUser[user.Id] = user
	return nil
}

func UpdateUser(user *User) error {
	if err := mysql.DB().Model(&User{}).Debug().Where(user.Id).Updates(user).Error; err != nil {
		return fmt.Errorf("create user failed, err: %s", err)
	}
	return nil
}

func UpdateUserRefreshToken(user *User, refreshToken string) error {
	if err := mysql.DB().Model(&User{}).Debug().Where(user).Update("refresh_token", refreshToken).Error; err != nil {
		return fmt.Errorf("create user failed, err: %s", err)
	}
	user.RefreshToken = refreshToken

	return nil
}

func GetUserByUsername(username string) (*User, error) {
	if user, ok := usernameToUser[username]; ok {
		return user, nil
	}
	return nil, UserNotExist
}

func GetUserByPhone(phone string) (*User, error) {
	if user, ok := phoneToUser[phone]; ok {
		return user, nil
	}
	return nil, UserNotExist
}

func GetUserByUid(uid int64) (*User, error) {
	if user, ok := uidToUser[uid]; ok {
		return user, nil
	}
	return nil, UserNotExist
}
