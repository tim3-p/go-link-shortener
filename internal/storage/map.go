package storage

import (
	"errors"
	"log"
)

type MapRepository struct {
	urlBase   map[string]string
	userLinks map[string][]string
}

func NewMapRepository() *MapRepository {
	return &MapRepository{urlBase: make(map[string]string), userLinks: make(map[string][]string)}
}

func (r *MapRepository) Add(key, value, userID string) error {
	log.Printf("Add userID - %s", userID)

	log.Printf("Add userLinks before - %s", r.userLinks)
	r.urlBase[key] = value
	r.userLinks[userID] = append(r.userLinks[userID], key)
	log.Printf("Add userLinks after - %s", r.userLinks)
	return nil
}

func (r *MapRepository) Get(key, userID string) (string, error) {
	log.Printf("Get userID - %s", userID)
	userMap := r.userLinks[userID]
	log.Printf("Get userLinks - %s", r.userLinks)
	log.Printf("Get userMap - %s", userMap)

	log.Printf("Get key - %s", key)

	for _, arrValue := range userMap {
		if arrValue == key {
			return r.urlBase[arrValue], nil
		}
	}

	return "", errors.New("key not found")
}

func (r *MapRepository) GetUserURLs(userID string) (map[string]string, error) {
	userMap := r.userLinks[userID]
	result := make(map[string]string)

	for _, arrValue := range userMap {
		if value, ok := r.urlBase[arrValue]; ok {
			result[arrValue] = value
		}
	}

	return result, nil
}
