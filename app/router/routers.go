package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"alligator/app/api"
	"alligator/pkg/utils/cache"
	"alligator/pkg/utils/log"
)

func InitRouter() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogMethod: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			log.Infof("Method: %s, URI: %s, status: %v", values.Method, values.URI, values.Status)
			return nil
		},
	}))
	root := e.Group("/v1")
	{
		root.POST("/login", api.Login)
		root.DELETE("/logout", api.Logout)

		g := root.Group("/")
		g.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
			Validator: func(key string, c echo.Context) (bool, error) {
				if _, err := cache.Get(key); err != nil {
					return false, nil
				}
				return true, nil
			},
		}))
		{
			g.GET("currentUser", api.CurrentUser)
			g.GET("users", api.GetUsers)
			g.POST("user", api.AddUser)
			g.POST("user/:id", api.UpdateUser)
			g.DELETE("user/:id", api.DeleteUser)
		}
	}

	return e
}
