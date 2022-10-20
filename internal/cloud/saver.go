package cloud

import (
	"context"

	"github.com/baibikov/tensile-cloud/internal/cloud/types"
)

type SaverMeta interface {
	SaveMeta(ctx context.Context, file types.File) (types.File, error)
}

type RemoverMeta interface {
	RemoveMeta(ctx context.Context, id string) error
}

type SaverFile interface {
	SaveFile(ctx context.Context, file types.File) error
}

type RemoverFile interface {
	RemoveFile(ctx context.Context, name string) error
}

type SaverRemoverMeta interface {
	SaverMeta
	RemoverMeta
}

type SaverRemoverFile interface {
	SaverFile
	RemoverFile
}

type Saver interface {
	SaverRemoverMeta
	SaverRemoverFile
}

type saverMetaCombine struct {
	SaverRemoverMeta
}

type saverFileCombine struct {
	SaverRemoverFile
}

type saverCombine struct {
	saverMetaCombine
	saverFileCombine
}

func NewSaver(meta SaverRemoverMeta, file SaverRemoverFile) Saver {
	return &saverCombine{
		saverMetaCombine: saverMetaCombine{
			SaverRemoverMeta: meta,
		},

		saverFileCombine: saverFileCombine{
			SaverRemoverFile: file,
		},
	}
}
