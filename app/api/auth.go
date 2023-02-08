package api

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"alligator/pkg/model"
)

type LoginUser struct {
	Username string `json:"username" form:"username" query:"username"`
	Password string `json:"password" form:"password"`
}

type LoginRes struct {
	Status string `json:"status"`
	Token  string `json:"token"`
}

func Login(c echo.Context) (err error) {
	u := new(LoginUser)
	if err = c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	user, _ := model.GetUser(u.Username)

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password)); err != nil {
		return c.JSON(http.StatusOK, &LoginRes{Status: "error", Token: ""})
	}

	token, err := model.GenerateJWT(u.Username)
	if err != nil {
		return ErrHandler(c, err)
	}

	return c.JSON(http.StatusOK, &LoginRes{Status: "ok", Token: token})
}

func Logout(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if token != "" {
		err := model.DeleteToken(strings.Replace(token, "bearer ", "", 1))
		if err != nil {
			return ErrHandler(c, err)
		}
	}
	return c.JSON(http.StatusNoContent, nil)
}
