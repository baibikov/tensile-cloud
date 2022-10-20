package cloud

import (
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/suite"

	"github.com/baibikov/tensile-cloud/internal/cloud/mocks"
)

// TODO: нужно много четких тестов.

//go:generate minimock -i Saver -o ./mocks/saver.go

type FileSuite struct {
	suite.Suite

	savermock  *mocks.SaverMock
	foldermock *mocks.FolderRepositoryMock
	folder     *File
}

func (f *FileSuite) SetupTest() {
	f.savermock = mocks.NewSaverMock(minimock.NewController(f.T()))
	f.foldermock = mocks.NewFolderRepositoryMock(minimock.NewController(f.T()))
	f.folder = &File{
		saver:      f.savermock,
		folderrepo: f.foldermock,
	}
}

func (f *FileSuite) TestUpload() {

}
