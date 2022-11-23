package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// The import below will automatically find and load the .env file.  Refer to
	// https://github.com/joho/godotenv/blob/main/autoload/autoload.go for details.
	_ "github.com/joho/godotenv/autoload"

	"github.com/NathanBak/cfgbuild"
	"github.com/NathanBak/picserv/internal/dbpics"
	"github.com/NathanBak/picserv/internal/server"
)

type Config struct {
	DbPicsConfig dbpics.Config `envvar:">,prefix=DROPBOX_"`
	ServerConfig server.Config `envvar:">"`
}

func main() {
	cfg, err := cfgbuild.NewConfig[*Config]()
	if err != nil {
		log.Fatal(err)
	}

	dbpics, err := dbpics.New(cfg.DbPicsConfig)
	if err != nil {
		log.Fatal(err)
	}

	cfg.ServerConfig.Picker = dbpics

	s, err := server.New(cfg.ServerConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Start server running on separate thread
	go func() {
		err = s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// wait for signal and then shutdown cleanly
	quitchan := make(chan os.Signal, 1)
	signal.Notify(quitchan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quitchan
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = s.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
