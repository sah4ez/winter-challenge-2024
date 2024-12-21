package main

/**
 * Grow and multiply your organisms to end up larger than your opponent.
 **/

func main() {
	game := NewGame()
	for {
		state := game.State()

		state.ScanStocks()
		state.ScanReqActions()

		state.Debug()
		state.DoAction()
	}
}
