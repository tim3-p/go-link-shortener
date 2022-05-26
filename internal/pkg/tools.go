package pkg

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

func HashURL(url []byte) string {
	hash := md5.Sum([]byte(url))
	str := hex.EncodeToString(hash[:])
	return string(str[1:6])
}

func CheckDBError(err error) (int, error) {
	var pgError *pgconn.PgError
	if errors.As(err, &pgError) && pgError.Code == pgerrcode.UniqueViolation {
		return http.StatusConflict, nil
	} else if err != nil {
		return http.StatusBadRequest, err
	}
	return http.StatusCreated, nil
}
