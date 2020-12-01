package store

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Store struct {
	config         *Config
	db             *sql.DB
	userRepository *UserRepository
}

func New(config *Config) *Store {
	return &Store{
		config: config,
	}
}

func (store *Store) Open() error {
	db, err := sql.Open("postgres", store.config.DatabaseUrl)

	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	return nil
}

func (store *Store) Close() error {

	store.db.Close()

	return nil
}

func (store *Store) User() *UserRepository {
	if store.userRepository != nil {
		return store.userRepository
	}

	store.userRepository = &UserRepository{
		store: store,
	}

	return store.userRepository
}
