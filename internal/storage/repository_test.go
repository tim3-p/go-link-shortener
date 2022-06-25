package storage

import (
	"testing"

	"github.com/google/uuid"
)

func BenchmarkGet(b *testing.B) {
	mapBase := NewMapRepository()
	arrayBase := NewArrayRepository()
	userID := uuid.NewString()
	URLs := map[string]string{
		"fsklfkdf": "http://aaaaa/index.php",
		"werwcxvx": "http://bbbbb/index.php",
		"lflsppsq": "http://ccccc/index.php",
		"rmakrwmv": "http://ddddd/index.php",
		"lpsbesow": "http://eeeee/index.php",
	}

	for arrKey, arrValue := range URLs {
		mapBase.Add(arrKey, arrValue, userID)
		arrayBase.Add(arrKey, arrValue, userID)
	}

	b.ReportAllocs()
	b.ResetTimer()

	b.Run("Map base ...", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			mapBase.Get("rmakrwmv", userID)
		}
	})

	b.Run("Array base ...", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			arrayBase.Get("rmakrwmv", userID)
		}
	})
}

func BenchmarkGetUserURLs(b *testing.B) {
	mapBase := NewMapRepository()
	arrayBase := NewArrayRepository()
	userID1 := uuid.NewString()
	userID2 := uuid.NewString()
	URLs1 := map[string]string{
		"fsklfkdf": "http://aaaaa/index.php",
		"werwcxvx": "http://bbbbb/index.php",
		"lflsppsq": "http://ccccc/index.php",
		"rmakrwmv": "http://ddddd/index.php",
		"lpsbesow": "http://eeeee/index.php",
	}

	URLs2 := map[string]string{
		"fsklfkdf": "http://aaaaa/index.php",
		"werwcxvx": "http://bbbbb/index.php",
		"lflsppsq": "http://ccccc/index.php",
		"rmakrwmv": "http://ddddd/index.php",
		"lpsbesow": "http://eeeee/index.php",
	}

	for arrKey, arrValue := range URLs1 {
		mapBase.Add(arrKey, arrValue, userID1)
		arrayBase.Add(arrKey, arrValue, userID1)
	}

	for arrKey, arrValue := range URLs2 {
		mapBase.Add(arrKey, arrValue, userID2)
		arrayBase.Add(arrKey, arrValue, userID2)
	}

	b.ReportAllocs()
	b.ResetTimer()

	b.Run("Map base ...", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			mapBase.GetUserURLs(userID1)
		}
	})

	b.Run("Array base ...", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			arrayBase.GetUserURLs(userID2)
		}
	})
}
