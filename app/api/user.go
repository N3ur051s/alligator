package api

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
	"golang.org/x/crypto/bcrypt"

	"alligator/pkg/model"
	"alligator/pkg/utils/cache"
)

type UserDTO struct {
	Name     string `json:"name" form:"name" query:"name"`
	Password string `json:"password" form:"password"`
	Email    string `json:"email" from:"email"`
	Phone    string `json:"phone" from:"phone"`
	IsAdmin  bool   `json:"isAdmin" form:"isAdmin"`
	Access   string `json:"access" from:"access"`
}

func CurrentUser(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	name, err := cache.Get(strings.Replace(token, "bearer ", "", 1))
	if err != nil {
		return ErrHandler(c, err)
	}
	user, err := model.GetUser(name)

	if err != nil {
		return ErrHandler(c, err)
	}
	return c.JSON(http.StatusOK, OkMes(user))
}

func GetUsers(c echo.Context) error {
	data := model.GetUserList(c, c.QueryParam("name"))

	return c.JSON(http.StatusOK, data)
}

func AddUser(c echo.Context) error {
	u := new(UserDTO)
	if err := c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	curd := model.NewCurd(&model.User{})

	pwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return ErrHandler(c, err)
	}
	u.Password = string(pwd)

	user := model.User{
		Name:     u.Name,
		Password: u.Password,
		IsAdmin:  false,
	}

	err = curd.Add(&user)

	if err != nil {
		return ErrHandler(c, err)
	}

	return c.JSON(http.StatusOK, user)
}

func UpdateUser(c echo.Context) error {
	userId := cast.ToInt(c.Param("id"))

	u := new(UserDTO)
	if err := c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	curd := model.NewCurd(&model.User{})

	var user, update model.User

	err := curd.First(&user, userId)

	if err != nil {
		return ErrHandler(c, err)
	}
	update.Name = u.Name

	// encrypt password
	if u.Password != "" {
		var pwd []byte
		pwd, err = bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return ErrHandler(c, err)
		}
		update.Password = string(pwd)
	}

	err = curd.Edit(&user, &update)

	if err != nil {
		return ErrHandler(c, err)
	}

	return c.JSON(http.StatusOK, user)
}

func DeleteUser(c echo.Context) error {
	id := c.Param("id")

	if cast.ToInt(id) == 1 {
		return c.JSON(http.StatusNotAcceptable, "message: Prohibit deleting the default user")
	}

	curd := model.NewCurd(&model.User{})
	err := curd.Delete(&model.User{}, "id", id)
	if err != nil {
		return ErrHandler(c, err)
	}
	return c.JSON(http.StatusNoContent, nil)
}
