package cloud

type Repository struct {
	File   FileRepository
	Folder FolderRepository
	Saver  Saver
}

type Cloud struct {
	repository *Repository
}

func (c Cloud) Folder() *Folder {
	return &Folder{
		folderrepo: c.repository.Folder,
	}
}

func (c Cloud) File() *File {
	return &File{
		filerepo:   c.repository.File,
		saver:      c.repository.Saver,
		folderrepo: c.repository.Folder,
	}
}

func New(repository *Repository) *Cloud {
	return &Cloud{
		repository: repository,
	}
}
