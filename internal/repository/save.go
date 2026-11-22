package repository

import "fmt"

func (r *Repository) Save(fullAddr, shortAddr string) (err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.storage[shortAddr] ; ok {
		return fmt.Errorf("value %s already exists", shortAddr)
	} 
	r.storage[shortAddr] = fullAddr
	return nil
}