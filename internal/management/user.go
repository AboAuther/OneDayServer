package management

import "errors"

var (
	usernameToUser = make(map[string]*User)
	phoneToUser    = make(map[string]*User)
	uidToUser      = make(map[int64]*User)

	UserNotExist = errors.New("user not exist")
)

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
