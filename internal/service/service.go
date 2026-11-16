package service


type Repository interface {
	Save(fullAddr, shortAddr string) (err error)
	Get(shortAddr string) (fullAddr string, err error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}
