package main

/**
 * Grow and multiply your organisms to end up larger than your opponent.
 **/

func main() {
	game := NewGame()
	// step := 0
	for {
		state := game.State()

		state.ScanStocks()
		state.ScanReqActions()

		state.DoAction(game)
		// full := true
		// state.Debug(full)
		// DebugMsg("step: ", step)
		// step += 1
	}
}
