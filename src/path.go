package main

func (s *State) PathScore(from Position, to Position) (float64, bool) {
	fromEntity := s.get(from.X, from.Y)
	toEntity := s.get(to.X, to.Y)
	// DebugMsg(">>>", fromEntity, toEntity)

	path := s.PathFind(fromEntity, toEntity)
	// DebugMsg(">>>", path)
	if path == nil {
		return 100, false
	}

	return path.TotalCost(), len(path.Entities) > 0
}

//func (s *State) PathScore2(from Position, to Position, depth int, hash map[string]struct{}, score int) (int, bool) {
//	find := false
//	if depth == MaxDepthWalking {
//		return 0, find
//	}
//	if score == InitScore {
//		return InitScore, find
//	}
//	if hash == nil {
//		hash = make(map[string]struct{}, 0)
//	}
//	hash[from.ID()] = struct{}{}
//
//	depth += 1
//	dirs := from.GetRoseLocality()
//	for i, dir := range dirs {
//		if !s.InMatrix(dir) {
//			continue
//		}
//		if _, ok := hash[dir.ID()]; ok {
//			continue
//		}
//		DebugMsg(">>", score, depth, from.ToLog(), to.ToLog(), dir.ToLog())
//		e := s.getByPos(dir)
//		if e != nil && (e.IsMy() || e.IsOpponent() || e.IsWall()) {
//			continue
//		}
//		if dir.Parent != nil {
//			continue
//		}
//		dir.Parent = &from
//		dirs[i] = dir
//		if dir.Equal(to) {
//			DebugMsg("FIND>", score, depth, from.ToLog(), to.ToLog(), dir.ToLog())
//			find = true
//			return score + 10, find
//		}
//	}
//
//	score += 10
//	for _, dir := range dirs {
//		if !s.InMatrix(dir) {
//			continue
//		}
//		if _, ok := hash[dir.ID()]; ok {
//			continue
//		}
//		// hash[dir.ID()] = struct{}{}
//		DebugMsg(">>", score, depth, from.ToLog(), to.ToLog(), dir.ToLog())
//		e := s.getByPos(dir)
//		if e != nil && (e.IsMy() || e.IsOpponent() || e.IsWall()) {
//			continue
//		}
//		score, find = s.PathScore(dir, to, depth, hash, score)
//		if find {
//			break
//		}
//	}
//	return score, find
//}
