package storage

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type DBRepository struct {
	connection *pgx.Conn
}

func NewDBRepository(pgConnection *pgx.Conn) (*DBRepository, error) {
	sql := `create table if not exists urls_base (
		short_url text not null primary key,
		original_url text,
		user_id      text		
	); 
	create unique index if not exists original_url_constrain on urls_base(original_url);`

	_, err := pgConnection.Exec(context.Background(), sql)
	if err != nil {
		return nil, err
	}
	return &DBRepository{connection: pgConnection}, nil
}

func (r *DBRepository) Add(key, value, userID string) error {
	sql := `insert into urls_base (short_url, original_url, user_id) values ($1, $2, $3)`
	_, err := r.connection.Exec(context.Background(), sql, key, value, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *DBRepository) Get(key, userID string) (string, error) {
	sql := `select original_url from urls_base where short_url = $1`
	row := r.connection.QueryRow(context.Background(), sql, key)
	var value string
	err := row.Scan(&value)
	if err != nil {
		return "", err
	}
	return value, nil
}

func (r *DBRepository) GetUserURLs(userID string) (map[string]string, error) {
	result := make(map[string]string)
	sql := `select short_url, original_url from urls_base where user_id = $1`
	rows, err := r.connection.Query(context.Background(), sql, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var shortURL, originalURL string
		err = rows.Scan(&shortURL, &originalURL)
		if err != nil {
			return nil, err
		}
		result[shortURL] = originalURL
	}
	return result, nil
}
