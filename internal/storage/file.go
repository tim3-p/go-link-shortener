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
	Key   string `json:"key"`
	Value string `json:"value"`
}

func NewFileRepository(fileStoragePath string) *FileRepository {
	return &FileRepository{fileStoragePath: fileStoragePath}
}

func (r *FileRepository) Add(key, value string) error {
	file, err := os.OpenFile(r.fileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(&FileRecord{Key: key, Value: value})
	if err != nil {
		return err
	}
	return nil
}

func (r *FileRepository) Get(key string) (string, error) {
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

		if record.Key == key {
			return record.Value, nil
		}
	}
	return "", errors.New("key not found")
}
