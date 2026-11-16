package repository

import "fmt"

func (r *Repository) Get(shortAddr string) (fullAddr string, err error) {
	fullAddr, ok := r.storage[shortAddr]
	if !ok {
		return "", fmt.Errorf("fullAddr not found for shortAddr: %s", shortAddr)
	}
	return fullAddr, nil
}
