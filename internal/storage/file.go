package storage

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

type FileRepository struct {
	fileStoragePath string
}

type FileRecord struct {
	UserID string `json:"user_id"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}

func NewFileRepository(fileStoragePath string) *FileRepository {
	return &FileRepository{fileStoragePath: fileStoragePath}
}

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
		} else if err != nil {
			result[record.Key] = record.Value
		}
	}
	return result, nil
}
