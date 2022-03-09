package app

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/google/uuid"
)

type cipherData struct {
	key    []byte
	nonce  []byte
	aesGCM cipher.AEAD
}

var cipherVal *cipherData

func cipherInit() error {
	if cipherVal == nil {
		key, err := generateRandom(2 * aes.BlockSize)
		if err != nil {
			return err
		}

		aesblock, err := aes.NewCipher(key)
		if err != nil {
			return err
		}

		aesgcm, err := cipher.NewGCM(aesblock)
		if err != nil {
			return err
		}

		nonce, err := generateRandom(aesgcm.NonceSize())
		if err != nil {
			return err
		}
		cipherVal = &cipherData{key: key, aesGCM: aesgcm, nonce: nonce}
	}
	return nil
}
func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func encrypt(userID string) (string, error) {
	if err := cipherInit(); err != nil {
		return "", err
	}
	encrypted := cipherVal.aesGCM.Seal(nil, cipherVal.nonce, []byte(userID), nil)
	return hex.EncodeToString(encrypted), nil
}

func decrypt(token string) (string, error) {
	if err := cipherInit(); err != nil {
		return "", err
	}
	b, err := hex.DecodeString(token)
	if err != nil {
		return "", err
	}
	userID, err := cipherVal.aesGCM.Open(nil, cipherVal.nonce, b, nil)
	if err != nil {
		return "", err
	}
	return string(userID), nil
}

func AuthHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := uuid.NewString()
		validAccessToken := false
		if c, err := r.Cookie("Token"); err == nil {
			if decrypted, err := decrypt(c.Value); err == nil {
				userID = decrypted
				validAccessToken = true
			}
		}
		if !validAccessToken {
			encrypted, err := encrypt(userID)
			if err != nil {
				http.Error(w, "Can not encrypt token", http.StatusInternalServerError)
				return
			}
			c := &http.Cookie{
				Name:  "Token",
				Value: encrypted,
				Path:  `/`,
			}
			http.SetCookie(w, c)
		}
		next.ServeHTTP(w, r)
	})
}
