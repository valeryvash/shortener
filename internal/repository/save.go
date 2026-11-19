package repository

func (r *Repository) Save(fullAddr, shortAddr string) (err error) {
	r.storage[shortAddr] = fullAddr
	return nil
}