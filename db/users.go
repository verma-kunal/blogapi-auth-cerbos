package db

import (
	"context"
	"errors"
)

type UserRecord struct {
	Password []byte
	Roles    []string
	Blocked  bool
}

var users = map[string]*UserRecord{
	"kunal": {
		Password: []byte(`$2y$10$s3QvUpMDYhdxO8LbPyiDou7KSTup.Hj9ip5ntB2h0NkW1fHIbYMm6`),
		Roles:    []string{"admin"},
		Blocked:  false,
	},
	"bella": {
		Password: []byte(`$2y$10$0V3N6CPkEozFKWhRgYSXJeXUEra2G7IYWr5BCSGwBSILRpLsfpVUm`),
		Roles:    []string{"user"},
		Blocked:  false,
	},
	"john": {
		Password: []byte(`$2y$10$RW1ItHGul1VXGZFs03YLFuwIvBijMv86uHq2pSHCkgvnvPHx10gj6`),
		Roles:    []string{"user"},
		Blocked:  false,
	},
}

// retrieve user info from the database
func FindUser(ctx context.Context, username string) (*UserRecord, error) {

	record, err := users[username]
	if !err {
		return nil, errors.New("record not found")
	}
	return record, nil
}
