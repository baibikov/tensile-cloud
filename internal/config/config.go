package config

import (
	"io/fs"
	"os"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"gopkg.in/yaml.v3"
)

type Config struct {
	API         API         `yaml:"api"`
	DB          DB          `yaml:"db"`
	FileStorage FileStorage `yaml:"file_storage"`
}

type (
	API struct {
		Host       string `yaml:"host"`
		UploadSize int64  `yaml:"upload_size"`
	}

	DB struct {
		Conn string `json:"conn"`
	}

	FileStorage struct {
		Dir string `yaml:"dir"`
	}
)

func New(path string) (config Config, err error) {
	file, err := os.OpenFile(path, os.O_RDONLY, fs.ModePerm)
	if err != nil {
		return config, errors.Wrapf(err, "open config by path %s", path)
	}
	defer func(err error) {
		multierr.AppendInto(&err, file.Close())
	}(err)

	return config, errors.Wrap(
		yaml.NewDecoder(file).Decode(&config),
		"decode config information",
	)
}
