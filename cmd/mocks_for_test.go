package cmd

import "github.com/UltimateSoftware/envctl/internal/db"

type memStore struct {
	env db.Environment
}

func (s *memStore) Create(e db.Environment) error {
	s.env = e

	return nil
}

func (s *memStore) Read() (db.Environment, error) {
	return s.env, nil
}

func (s *memStore) Delete() error {
	s.env = db.Environment{}
	return nil
}
