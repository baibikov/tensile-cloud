package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.uber.org/multierr"

	"github.com/baibikov/tensile-cloud/internal/config"
	"github.com/baibikov/tensile-cloud/internal/files"
	"github.com/baibikov/tensile-cloud/internal/swagger"
	"github.com/baibikov/tensile-cloud/pkg/httpserver"
	"github.com/baibikov/tensile-cloud/pkg/multistore"
)

func init() {
	logLevel := os.Getenv("log_level")
	if logLevel == "" {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
}

const (
	configPath = "./configs/main.yaml"
)

func runApp() (err error) {
	logrus.Info("init app context")
	ctx := context.Background()

	logrus.Info("init app notification")
	ctx, cancel := signal.NotifyContext(
		ctx,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
	)
	defer cancel()

	logrus.Info("init app config")
	cconfig, err := config.New(configPath)
	if err != nil {
		return errors.Wrap(err, "init app config")
	}

	logrus.Info("init app store")
	store, err := multistore.New(&multistore.Config{
		Conn: cconfig.DB.Conn,
	})
	if err != nil {
		return errors.Wrap(err, "init app multistore")
	}
	defer func() {
		multierr.AppendInto(&err, store.Close())
	}()

	cloudhandler, err := initCloud(cconfig, store)
	if err != nil {
		return err
	}

	logrus.Info("init app server")
	server := httpserver.New()

	server.AddWithFullPath("/api/cloud", cloudhandler)
	server.AddWithFullPath("/files", files.New("./.data"))
	server.AddWithFullPath("/swagger-ui", swagger.New("./static/swagger-ui"))

	go func() {
		logrus.Infof("start listen app server by host: %s", cconfig.API.Host)
		if serr := server.Serve(cconfig.API.Host); serr != nil {
			multierr.AppendInto(&err, serr)
			cancel()
		}
	}()

	<-ctx.Done()
	logrus.Info("app shutdown")
	time.Sleep(time.Second)

	return nil
}
