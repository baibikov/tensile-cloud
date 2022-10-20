package transaction

import "go.uber.org/multierr"

type TX struct {
	commit   CommitFunc
	rollback RollbackFunc
}

func New(commit CommitFunc, rollback RollbackFunc) *TX {
	return &TX{
		commit:   commit,
		rollback: rollback,
	}
}

type (
	CommitFunc   func() error
	RollbackFunc func() error
)

func (t *TX) Run() (err error) {
	defer func() {
		if err != nil {
			multierr.AppendInto(&err, t.rollback())
		}
	}()

	return t.commit()
}
