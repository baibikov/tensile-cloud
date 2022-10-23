package cloud

import (
	"context"
	"testing"
	"time"

	"github.com/go-openapi/swag"
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
	defer f.foldermock.MinimockGetByParentDone()
	ctx := context.Background()
	mocktime := time.Now()

	f.foldermock.GetByParentMock.When(
		ctx,
		nil,
		types.NewSort(swag.String("name"), nil),
	).Then([]*types.Folder{
		{
			ID:        "1",
			Name:      "folder_test_1",
			CreatedAt: mocktime,
			UpdatedAt: mocktime,
		},
		{
			ID:        "2",
			Name:      "folder_test_2",
			CreatedAt: mocktime,
			UpdatedAt: mocktime,
		},
		{
			ID:        "3",
			Name:      "folder_test_3",
			CreatedAt: mocktime,
			UpdatedAt: mocktime,
		},
	}, nil)

	f.foldermock.GetByParentMock.When(
		ctx,
		nil,
		types.NewSort(swag.String("name"), swag.String("desc")),
	).Then([]*types.Folder{
		{
			ID:        "3",
			Name:      "folder_test_3",
			CreatedAt: mocktime,
			UpdatedAt: mocktime,
		},
		{
			ID:        "2",
			Name:      "folder_test_2",
			CreatedAt: mocktime,
			UpdatedAt: mocktime,
		},
		{
			ID:        "1",
			Name:      "folder_test_1",
			CreatedAt: mocktime,
			UpdatedAt: mocktime,
		},
	}, nil)

	f.foldermock.GetByParentMock.When(ctx,
		swag.String("generated_uuid"),
		types.NewSort(swag.String("name"), nil),
	).Then([]*types.Folder{
		{
			ID:        "4",
			ParentID:  swag.String("generated_uuid"),
			Name:      "folder_test_4",
			CreatedAt: mocktime,
			UpdatedAt: mocktime,
		},
		{
			ID:        "5",
			ParentID:  swag.String("generated_uuid"),
			Name:      "folder_test_5",
			CreatedAt: mocktime,
			UpdatedAt: mocktime,
		},
	}, nil)

	tests := []struct {
		name string
		args struct {
			ctx      context.Context
			parentID *string
			sort     types.Sort
		}
		want []*types.Folder
	}{
		{
			name: "null parent with sort asc",
			args: struct {
				ctx      context.Context
				parentID *string
				sort     types.Sort
			}{
				ctx:      ctx,
				parentID: nil,
				sort:     types.NewSort(swag.String("name"), nil),
			},
			want: []*types.Folder{
				{
					ID:        "1",
					Name:      "folder_test_1",
					CreatedAt: mocktime,
					UpdatedAt: mocktime,
				},
				{
					ID:        "2",
					Name:      "folder_test_2",
					CreatedAt: mocktime,
					UpdatedAt: mocktime,
				},
				{
					ID:        "3",
					Name:      "folder_test_3",
					CreatedAt: mocktime,
					UpdatedAt: mocktime,
				},
			},
		},
		{
			name: "null parent with sort desc",
			args: struct {
				ctx      context.Context
				parentID *string
				sort     types.Sort
			}{
				ctx:      ctx,
				parentID: nil,
				sort:     types.NewSort(swag.String("name"), swag.String("desc")),
			},
			want: []*types.Folder{
				{
					ID:        "3",
					Name:      "folder_test_3",
					CreatedAt: mocktime,
					UpdatedAt: mocktime,
				},
				{
					ID:        "2",
					Name:      "folder_test_2",
					CreatedAt: mocktime,
					UpdatedAt: mocktime,
				},
				{
					ID:        "1",
					Name:      "folder_test_1",
					CreatedAt: mocktime,
					UpdatedAt: mocktime,
				},
			},
		},
		{
			name: "has parent with sort asc",
			args: struct {
				ctx      context.Context
				parentID *string
				sort     types.Sort
			}{
				ctx:      ctx,
				parentID: swag.String("generated_uuid"),
				sort:     types.NewSort(swag.String("name"), nil),
			},
			want: []*types.Folder{
				{
					ID:        "4",
					ParentID:  swag.String("generated_uuid"),
					Name:      "folder_test_4",
					CreatedAt: mocktime,
					UpdatedAt: mocktime,
				},
				{
					ID:        "5",
					ParentID:  swag.String("generated_uuid"),
					Name:      "folder_test_5",
					CreatedAt: mocktime,
					UpdatedAt: mocktime,
				},
			},
		},
	}

	for _, tt := range tests {
		f.Run(tt.name, func() {
			actual, err := f.folder.Find(tt.args.ctx, tt.args.parentID, tt.args.sort)
			f.Require().NoError(err)
			f.Require().Equal(tt.want, actual)
		})
	}
}

func TestFolderSuite(t *testing.T) {
	suite.Run(t, new(FolderSuite))
}
