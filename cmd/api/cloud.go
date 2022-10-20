package main

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/baibikov/tensile-cloud/internal/cloud"
	"github.com/baibikov/tensile-cloud/internal/cloud/repository/filestorage"
	"github.com/baibikov/tensile-cloud/internal/cloud/repository/postgres"
	"github.com/baibikov/tensile-cloud/internal/cloud/rest"
	"github.com/baibikov/tensile-cloud/internal/config"
	"github.com/baibikov/tensile-cloud/pkg/httpserver"
	"github.com/baibikov/tensile-cloud/pkg/multistore"
)

func initCloud(config config.Config, store multistore.Store) (httpserver.Handler, error) {
	pgfilerepo := postgres.NewFile(store.DB())
	fsfilerepo := filestorage.NewFile(config.FileStorage.Dir)

	logrus.Info("init app cloud repository")
	cloudrepo := &cloud.Repository{
		File:   newCombine(pgfilerepo, fsfilerepo),
		Folder: postgres.NewFolder(store.DB()),
		Saver:  cloud.NewSaver(pgfilerepo, fsfilerepo),
	}

	logrus.Info("init app cloud usecase")
	ccloud := cloud.New(cloudrepo)

	restConfig := &rest.Config{
		UploadSize: config.API.UploadSize,
	}

	restUseCase := &rest.UseCase{
		Folder: ccloud.Folder(),
		File:   ccloud.File(),
	}

	logrus.Info("init app cloud rest handlers")
	handler, err := rest.New(restConfig, restUseCase)
	return handler, errors.Wrap(err, "init app cloud rest handlers")
}

type pgComponent struct {
	*postgres.File
}

type fsComponent struct {
	*filestorage.File
}

type combine struct {
	*pgComponent
	*fsComponent
}

func newCombine(pg *postgres.File, fs *filestorage.File) *combine {
	return &combine{
		pgComponent: &pgComponent{
			File: pg,
		},
		fsComponent: &fsComponent{
			File: fs,
		},
	}
}
