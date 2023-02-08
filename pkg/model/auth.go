package model

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	"alligator/pkg/utils/cache"
)

type User struct {
	Model

	Name     string `json:"name"`
	Password string `json:"-"`
	IsAdmin  bool   `json:"isAdmin" gorm:"column:isAdmin"`
}

type AuthToken struct {
	Token string `json:"token"`
}

type JWTClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

func GetUser(name string) (user User, err error) {
	err = db.Where("name = ?", name).First(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, err
}

func GetUserList(c echo.Context, username interface{}) (data DataList) {
	var total int64
	db.Model(&User{}).Count(&total)
	var users []User

	result := db.Model(&User{}).Scopes(orderAndPaginate(c))

	if username != "" {
		result = result.Where("name LIKE ?", "%"+username.(string)+"%")
	}

	result.Find(&users)

	data = GetListWithPagination(&users, c, total)
	return
}

func DeleteToken(token string) error {
	return cache.Del(token)
}

func GenerateJWT(name string) (string, error) {
	claims := JWTClaims{
		Name: name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
	}
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := unsignedToken.SignedString([]byte("2E1CE615-BB15-44F5-B5BE-6B5DA3581D0F"))
	if err != nil {
		return "", err
	}

	if err = cache.Set(signedToken, name, time.Hour); err != nil {
		return "", err
	}

	return signedToken, err
}
