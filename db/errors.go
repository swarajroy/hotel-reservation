package db

type DBError struct {
	Err string
}

func NewResourceError(err string) error {
	return DBError{
		Err: err,
	}
}

func (e DBError) Error() string {
	return e.Err
}
