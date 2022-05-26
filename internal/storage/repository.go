package storage

var urlBase = make(map[string]string)

func Add(key, value string) {
	urlBase[key] = value
}

func Get(key string) (string, bool) {
	if value, ok := urlBase[key]; ok {
		return value, true
	} else {
		return "", false
	}
}
