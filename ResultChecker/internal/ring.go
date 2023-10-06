package internal

import "container/ring"

// RemoveFromRing -> take a ring and remove element with specified index
func RemoveFromRing(r *ring.Ring, index int) *ring.Ring {
	if index < 0 || index >= r.Len() {
		// Not valid index
		return r
	}

	newRing := ring.New(r.Len() - 1)
	counter := 0

	r.Do(func(p interface{}) {
		if counter != index {
			newRing.Value = p
			newRing = newRing.Next()
		}
		counter++
	})

	return newRing
}
