package storage

import (
	"errors"
	"log"
)

// InMemory base descriptor
type MapRepository struct {
	urlBase   map[string]string
	userLinks map[string][]string
}

// Constructor for InMemory storage
func NewMapRepository() *MapRepository {
	return &MapRepository{urlBase: make(map[string]string), userLinks: make(map[string][]string)}
}

// Add new short URL in map storage
func (r *MapRepository) Add(key, value, userID string) error {
	log.Printf("Add userID - %s", userID)
	log.Printf("Add userLinks before - %s", r.userLinks)
	r.urlBase[key] = value
	r.userLinks[userID] = append(r.userLinks[userID], key)
	log.Printf("Add userLinks after - %s", r.userLinks)
	return nil
}

// Get origin URL by short URL from map storage
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
	/*
		if value, ok := r.urlBase[key]; ok {
			return value, nil
		} else {
			return "", errors.New("key not found")
		}
	*/
}

// Get URLs by user ID from map storage
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

// Delete URL for user ID from map storage
func (r *MapRepository) Delete(keys []string, userID string) error {
	return nil
}
