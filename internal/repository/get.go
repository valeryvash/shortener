package repository

import "fmt"

func (r *Repository) Get(shortAddr string) (fullAddr string, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	fullAddr, ok := r.storage[shortAddr]
	if !ok {
		return "", fmt.Errorf("fullAddr not found for shortAddr: %s", shortAddr)
	}
	return fullAddr, nil
}
