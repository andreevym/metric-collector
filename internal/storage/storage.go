package storage

type Storage interface {
	CreateAll(kvMap map[string]*Metric) error
	Create(key string, val string) error
	Read(key string) (string, error)
	Update(key string, val string) error
	Delete(key string) error
}

type Metric struct {
	Value    string
	IsExists bool
}
