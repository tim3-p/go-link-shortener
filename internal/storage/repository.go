package storage

type Repository interface {
	Add(key, value, userID string) error
	Get(key, userID string) (string, error)
	GetUserURLs(userID string) (map[string]string, error)
	Delete(keys []string, userID string) error
}
