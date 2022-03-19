package storage

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type DbRepository struct {
	connection *pgx.Conn
}

func NewDbRepository(pgConnection *pgx.Conn) (*DbRepository, error) {
	sql := `create table if not exists urls_base (
		short_url text not null primary key,
		original_url text,
		user_id      text		
	);`
	_, err := pgConnection.Exec(context.Background(), sql)
	if err != nil {
		return nil, err
	}
	return &DbRepository{connection: pgConnection}, nil
}

func (r *DbRepository) Add(key, value, userID string) error {
	sql := `insert into urls_base (short_url, original_url, user_id) values ($1, $2, $3)`
	_, err := r.connection.Exec(context.Background(), sql, key, value, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *DbRepository) Get(key, userID string) (string, error) {
	sql := `select original_url from urls_base where short_url = $1`
	row := r.connection.QueryRow(context.Background(), sql, key)
	var value string
	err := row.Scan(&value)
	if err != nil {
		return "", err
	}
	return value, nil
}

func (r *DbRepository) GetUserURLs(userID string) (map[string]string, error) {
	result := make(map[string]string)
	sql := `select short_url, original_url from urls where user_id = $1`
	rows, err := r.connection.Query(context.Background(), sql, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var short_url, original_url string
		err = rows.Scan(&short_url, &original_url)
		if err != nil {
			return nil, err
		}
		result[short_url] = original_url
	}
	return result, nil
}
