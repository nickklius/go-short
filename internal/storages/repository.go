package storages

type Repository interface {
	Create(longURL string) (string, error)
	Read(shortURL string) (string, error)
	GetAll() *map[string]string
}

func CreateShortURL(r Repository, longURL string) (string, error) {
	return r.Create(longURL)
}

func RetrieveURL(r Repository, shortURL string) (string, error) {
	return r.Read(shortURL)
}
