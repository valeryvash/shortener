package service

import "fmt"

func (s *Service) GetFullAddr(shortAddr string) (fullAddr string, err error) {
	fullAddr, err = s.repo.Get(shortAddr)
	if err != nil {
		return "", fmt.Errorf("failed to get: %w", err)
	}

	return fullAddr, nil
}