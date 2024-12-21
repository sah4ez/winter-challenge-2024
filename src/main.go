package main

import "fmt"

/**
 * Grow and multiply your organisms to end up larger than your opponent.
 **/

func main() {
	// width: columns in the game grid
	// height: rows in the game grid
	var width, height int
	fmt.Scan(&width, &height)

	for {
		var entityCount int
		fmt.Scan(&entityCount)

		for i := 0; i < entityCount; i++ {
			// y: grid coordinate
			// _type: WALL, ROOT, BASIC, TENTACLE, HARVESTER, SPORER, A, B, C, D
			// owner: 1 if your organ, 0 if enemy organ, -1 if neither
			// organId: id of this entity if it's an organ, 0 otherwise
			// organDir: N,E,S,W or X if not an organ
			var x, y int
			var _type string
			var owner, organId int
			var organDir string
			var organParentId, organRootId int
			fmt.Scan(&x, &y, &_type, &owner, &organId, &organDir, &organParentId, &organRootId)
		}
		// myD: your protein stock
		var myA, myB, myC, myD int
		fmt.Scan(&myA, &myB, &myC, &myD)

		// oppD: opponent's protein stock
		var oppA, oppB, oppC, oppD int
		fmt.Scan(&oppA, &oppB, &oppC, &oppD)

		// requiredActionsCount: your number of organisms, output an action for each one in any order
		var requiredActionsCount int
		fmt.Scan(&requiredActionsCount)

		for i := 0; i < requiredActionsCount; i++ {

			// fmt.Fprintln(os.Stderr, "Debug messages...")
			fmt.Println("WAIT") // Write action to stdout
		}
	}
}
