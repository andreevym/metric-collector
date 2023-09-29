package mem

type Storage struct {
	m map[string][]string
}

func NewStorage() Storage {
	return Storage{
		map[string][]string{},
	}
}

func (s Storage) Create(key string, val string) error {
	values := s.m[key]
	if values == nil {
		values = make([]string, 0)
	}
	s.m[key] = append(values, val)
	return nil
}

func (s Storage) Read(key string) ([]string, error) {
	return s.m[key], nil
}

func (s Storage) Update(key string, val []string) error {
	s.m[key] = val
	return nil
}

func (s Storage) Delete(key string) error {
	delete(s.m, key)
	return nil
}
