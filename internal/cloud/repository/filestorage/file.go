package filestorage

import (
	"context"
	"io"
	"os"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/baibikov/tensile-cloud/internal/cloud/types"
)

type File struct {
	dir string
}

func NewFile(dir string) *File {
	return &File{dir: dir}
}

func (f File) SaveFile(_ context.Context, file types.File) (err error) {
	filename := file.ID
	if file.Type != "" {
		filename += "." + file.Type
	}

	src, err := os.Create(f.dir + "/" + filename)
	if err != nil {
		return errors.Wrapf(err, "create src file by name - %s", filename)
	}
	defer func() {
		multierr.AppendInto(&err, src.Close())
	}()

	_, err = io.Copy(src, file.Body)
	return errors.Wrap(err, "copy from dest to src")
}

func (f File) RemoveFile(_ context.Context, name string) error {
	return os.Remove(f.dir + "/" + name)
}

func (f File) Open(_ context.Context, name string) (io.ReadCloser, error) {
	return os.Open(f.dir + "/" + name)
}

func (f File) Copy(_ context.Context, dst, src string) (err error) {
	srcf, err := os.Open(f.dir + "/" + src)
	if err != nil {
		return err
	}
	defer func() {
		multierr.AppendInto(&err, srcf.Close())
	}()

	dstf, err := os.Create(f.dir + "/" + dst)
	if err != nil {
		return err
	}
	defer func() {
		multierr.AppendInto(&err, dstf.Close())
	}()

	_, err = io.Copy(dstf, srcf)
	return err
}
