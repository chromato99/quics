package sync

import "github.com/dgraph-io/badger/v3"

type Repository struct {
	DB *badger.DB
}

type RepositoryInterface interface {
}

func NewSyncRepository(db *badger.DB) *Repository {
	return &Repository{DB: db}
}
