package main

func (s *State) walk(x, y int, fn func(e *Entity) bool) {
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	if x >= len(s.matrix) || len(s.matrix) == 0 {
		return
	}
	if y >= len(s.matrix[0]) {
		return
	}

	for _, k := range s.matrix[x:] {
		if k == nil {
			continue
		}
		for _, n := range k[y:] {
			if n != nil && n.IsMy() && s.FreeEntites(n) {
				if !fn(n) {
					break
				}
			}
		}
	}
}
