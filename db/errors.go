package db

type DBError struct {
	Err string
}

func NewDBError(err string) DBError {
	return DBError{
		Err: err,
	}
}

func (e DBError) Error() string {
	return e.Err
}
