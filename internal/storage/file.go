package storage

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

// FileRepository File base descriptor
type FileRepository struct {
	fileStoragePath string
}

// File record structure
type FileRecord struct {
	UserID string `json:"user_id"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}

// Constructor for file storage
func NewFileRepository(fileStoragePath string) *FileRepository {
	return &FileRepository{fileStoragePath: fileStoragePath}
}

// Add new short URL in file storage
func (r *FileRepository) Add(key, value, userID string) error {
	file, err := os.OpenFile(r.fileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(&FileRecord{UserID: userID, Key: key, Value: value})
	if err != nil {
		return err
	}
	return nil
}

// Get origin URL by short URL from file storage
func (r *FileRepository) Get(key, userID string) (string, error) {
	file, err := os.OpenFile(r.fileStoragePath, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return "", err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	for {
		record := &FileRecord{}
		if err := decoder.Decode(&record); err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}

		if record.Key == key && record.UserID == userID {
			return record.Value, nil
		}
	}
	return "", errors.New("key not found")
}

// Get URLs by user ID from file storage
func (r *FileRepository) GetUserURLs(userID string) (map[string]string, error) {
	result := make(map[string]string)

	file, err := os.OpenFile(r.fileStoragePath, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return result, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	for {
		record := &FileRecord{}
		if err := decoder.Decode(&record); err == io.EOF {
			break
		} else if err != nil && record.UserID == userID {
			result[record.Key] = record.Value
		}
	}
	return result, nil
}

// Delete URL for user ID from file storage
func (r *FileRepository) Delete(keys []string, userID string) error {
	return nil
}

// Returs stats from file storage
func (r *FileRepository) GetStats() (int, int, error) {
	return 0, 0, nil
}
