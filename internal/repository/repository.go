package repository

type Storage interface {
	Create(key string, val string) error
	Read(key string) ([]string, error)
	Update(key string, val []string) error
	Delete(key string) error
}
