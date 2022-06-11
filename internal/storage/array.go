package storage

import (
	"errors"
)

type ArrayRecord struct {
	UserID string `json:"user_id"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}

// Array base descriptor
type ArrayRepository struct {
	urlBase []ArrayRecord
}

// Constructor for Array storage
func NewArrayRepository() *ArrayRepository {
	return &ArrayRepository{urlBase: make([]ArrayRecord, 0)}
}

// Add new short URL in array storage
func (r *ArrayRepository) Add(key, value, userID string) error {
	row := ArrayRecord{
		UserID: userID,
		Key:    key,
		Value:  value,
	}
	r.urlBase = append(r.urlBase, row)
	return nil
}

// Get origin URL by short URL from array storage
func (r *ArrayRepository) Get(key, userID string) (string, error) {
	for _, arrValue := range r.urlBase {
		if arrValue.Key == key && arrValue.UserID == userID {
			return arrValue.Value, nil
		}
	}
	return "", errors.New("key not found")
}

// Get URLs by user ID from array storage
func (r *ArrayRepository) GetUserURLs(userID string) (map[string]string, error) {
	result := make(map[string]string)
	for _, arrValue := range r.urlBase {
		if arrValue.UserID == userID {
			result[arrValue.Key] = arrValue.Value
		}
	}
	return result, nil
}

// Delete URL for user ID from array storage
func (r *ArrayRepository) Delete(keys []string, userID string) error {
	return nil
}
