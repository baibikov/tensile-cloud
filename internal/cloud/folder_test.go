package cloud

import (
	"context"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/suite"

	"github.com/baibikov/tensile-cloud/internal/cloud/mocks"
	"github.com/baibikov/tensile-cloud/internal/cloud/types"
)

//go:generate minimock -i FolderRepository -o ./mocks/folder_repository.go

type FolderSuite struct {
	suite.Suite

	foldermock *mocks.FolderRepositoryMock
	folder     *Folder
}

func (f *FolderSuite) SetupTest() {
	f.foldermock = mocks.NewFolderRepositoryMock(minimock.NewController(f.T()))
	f.folder = &Folder{
		folderrepo: f.foldermock,
	}
}

func (f *FolderSuite) TestFind() {
	mockTime := time.Now()
	mockCtx := context.Background()

	mockFolders1 := []*types.Folder{
		{
			ID:        "1",
			Name:      "testname",
			CreatedAt: mockTime,
			UpdatedAt: mockTime,
		},
	}
	mockFolders2 := []*types.Folder{
		{
			ID:        "2",
			Name:      "testname",
			CreatedAt: mockTime,
			UpdatedAt: mockTime,
		},
	}

	f.foldermock.GetByParentMock.Set(func(ctx context.Context, parentID *string) (fpa1 []*types.Folder, err error) {
		if parentID == nil {
			return mockFolders1, nil
		}

		return mockFolders2, nil
	})

	ff, err := f.folder.Find(mockCtx, nil)
	f.NoError(err)

	f.Equal(ff, mockFolders1)
}

func TestFolderSuite(t *testing.T) {
	suite.Run(t, new(FolderSuite))
}
