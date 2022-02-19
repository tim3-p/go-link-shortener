package storage

import "errors"

var urlBase = make(map[string]string)

type Repository interface {
	Add(key, value string) error
	Get(key string) (string, error)
}

func Add(key, value string) error {
	urlBase[key] = value
	return nil
}

func Get(key string) (string, error) {
	if value, ok := urlBase[key]; ok {
		return value, nil
	} else {
		return "", errors.New("key not found")
	}
}
