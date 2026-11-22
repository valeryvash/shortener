package service

import (
	"crypto/sha256"
	"fmt"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var runes = []rune(alphabet)

func (s *Service) CreateShortAddr(fullAddr string) (shortAddr string, err error) {
	hash := sha256.Sum256([]byte(fullAddr))

	result := make([]rune, 8)
	for i := range result {
		result[i] = runes[int(hash[i])%len(runes)]
	}

	shortAddr = string(result)

	err = s.repo.Save(fullAddr, shortAddr)
	if err != nil {
		return "", fmt.Errorf("failed to persist: %w", err)
	}

	return shortAddr, nil
}
