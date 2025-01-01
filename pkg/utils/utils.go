package utils

import (
	"chat-application/users/config"
	"chat-application/users/pkg/jwt"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Utils struct{}

func NewUtils() *Utils {
	return &Utils{}
}

func (u *Utils) HashPassword(password string) (string, error) {
	hashMult := 10
	byte, err := bcrypt.GenerateFromPassword([]byte(password), hashMult)

	if err != nil {
		return "", err
	}

	return string(byte), nil
}

func (u *Utils) ComparePassword(password, inputPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(inputPassword))
}

func (u *Utils) GenerateLoginToken(id uint, username string) (string, string, error) {
	env, _ := config.NewEnvDev()
	accessJwtHelper := jwt.NewJWTHelper(env.JWT_SECRET, 15*time.Minute)
	refreshJwtHelper := jwt.NewJWTHelper(env.JWT_SECRET, 24*time.Hour)

	access, err := accessJwtHelper.GenerateToken(u.UintToString(id), username)
	if err != nil {
		return "", "", err
	}

	refresh, err := refreshJwtHelper.GenerateToken(u.UintToString(id), username)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (u *Utils) UintToString(arg uint) string {
	return string(strconv.FormatUint(uint64(arg), 10))
}
