package storage

import "errors"

type MapRepository struct {
	urlBase map[string]string
}

func NewMapRepository() *MapRepository {
	return &MapRepository{urlBase: make(map[string]string)}
}

func (r *MapRepository) Add(key, value string) error {
	r.urlBase[key] = value
	return nil
}

func (r *MapRepository) Get(key string) (string, error) {
	if value, ok := r.urlBase[key]; ok {
		return value, nil
	} else {
		return "", errors.New("key not found")
	}
}

func (r *MapRepository) GetUserURLs() (map[string]string, error) {
	return r.urlBase, nil
}
