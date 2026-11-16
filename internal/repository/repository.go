package repository

type Repository struct {
	storage map[string]string
}

func NewRepository() *Repository {
	return &Repository{storage: make(map[string]string)}
}
