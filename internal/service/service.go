package service

//go:generate mockery --name=Repository --output=mocks --outpkg=mocks
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
