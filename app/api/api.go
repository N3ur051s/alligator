package api

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"alligator/pkg/utils/log"
)

type Res struct {
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
	ErrMes  string      `json:"errorMessage"`
}

func ErrMes(err string) *Res {
	return &Res{Success: false, Data: "", ErrMes: err}
}

func ErrHandler(c echo.Context, err error) error {
	log.Error(err)
	return c.JSON(http.StatusInternalServerError, ErrMes(err.Error()))
}

func OkMes(data interface{}) *Res {
	return &Res{Success: true, Data: data}
}
