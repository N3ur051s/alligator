package alligator

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"alligator/app/router"
	"alligator/pkg/config"
	"alligator/pkg/utils/log"
)

func Run(opts config.Options) {
	log.Infof("Alligator arguments %+v", opts)

	e := router.InitRouter()
	// Start server
	go func() {
		if err := e.Start(":" + strconv.Itoa(opts.HTTPListenPort)); err != nil && err != http.ErrServerClosed {
			log.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
