package main

import (
	_ "net/http/pprof"
	"os"

	"github.com/urfave/cli/v2"

	"alligator/app/alligator"
	"alligator/pkg/config"
	"alligator/pkg/model"
	"alligator/pkg/utils/cache"
	"alligator/pkg/utils/log"
)

func main() {
	// management.RegisterPassword()

	var config config.Options

	app := cli.NewApp()
	app.Usage = "Complete container management platform"
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:        "debug",
			Usage:       "Enable debug logs",
			Destination: &config.Debug,
		},
		&cli.IntFlag{
			Name:        "http-listen-port",
			Usage:       "HTTP listen port",
			Value:       8080,
			Destination: &config.HTTPListenPort,
		},
		&cli.StringFlag{
			Name:  "log-format",
			Usage: "Log formatter used (json, text, simple)",
			Value: "simple",
		},
		&cli.StringFlag{
			Name:        "audit-log-path",
			EnvVars:     []string{"AUDIT_LOG_PATH"},
			Value:       "/var/log/alligator/api-audit.log",
			Usage:       "Log path for Alligator Server API. Default path is /var/log/alligator/api-audit.log",
			Destination: &config.AuditLogPath,
		},
		&cli.StringFlag{
			Name:        "db-ip",
			EnvVars:     []string{"DB_IP"},
			Value:       "localhost",
			Usage:       "DataBase IP Address. Default localhost",
			Destination: &config.Db.Ip,
		},
		&cli.IntFlag{
			Name:        "db-port",
			EnvVars:     []string{"DB_PORT"},
			Value:       3306,
			Usage:       "DataBase Port. Default Port is 3306",
			Destination: &config.Db.Port,
		},
		&cli.StringFlag{
			Name:        "db-user",
			EnvVars:     []string{"DB_USER"},
			Value:       "root",
			Usage:       "User for login Database. Default User is root",
			Destination: &config.Db.User,
		},
		&cli.StringFlag{
			Name:        "db-passwd",
			EnvVars:     []string{"DB_PASSWD"},
			Usage:       "DataBase Passwd",
			Destination: &config.Db.Passwd,
		},
		&cli.StringFlag{
			Name:        "redis-addr",
			EnvVars:     []string{"REDIS_ADDR"},
			Value:       "127.0.0.1:6379",
			Usage:       "CACHE IP Address. Default 127.0.0.1:6379",
			Destination: &config.Cache.Addr,
		},
		&cli.StringFlag{
			Name:        "redis-passwd",
			EnvVars:     []string{"REDIS_PASSWD"},
			Usage:       "REDIS PASSWORD. ",
			Destination: &config.Cache.Passwd,
		},
		&cli.IntFlag{
			Name:        "redis-db",
			EnvVars:     []string{"REDIS_DB"},
			Value:       0,
			Usage:       "REDIS DB. Default 0",
			Destination: &config.Cache.DB,
		},
	}

	app.Action = func(c *cli.Context) error {
		log.Init(c, config)
		model.Init(config)
		cache.Init(config)
		alligator.Run(config)
		return nil
	}

	app.ExitErrHandler = func(c *cli.Context, err error) {
		log.Fatal(err)
	}

	app.Run(os.Args)
}
