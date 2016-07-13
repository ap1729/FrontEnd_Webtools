package perf

import (
	"FrontEnd_WebTools/model"
	"FrontEnd_WebTools/service"
	"fmt"
)

// Returns a list of ID's that identify interfering BaseStations for a user in
// the scenario with given frequency-reuse mode.
func intrStations(sc *model.Scenario, hexMap *service.HexMap, userID uint, p *Params) []uint {
	var bsIds []uint
	switch p.FrMode {

	case "FR1":
		// FR1 sends the list of BS in cells that match the operator(s) enabled
		bsIds = *new([]uint)
		for i := uint(0); i < uint(len(sc.BaseStations())); i++ {
			if p.OpEnableFlags[sc.GetStationByID(i).OwnerOp().ID()] == true {
				bsIds = append(bsIds, i)
			}
		}

	case "FR3":
		//FR3 sends a list of BS in FR3 cells that match the operator(s) enabled.
		//Get x,y locations of UE and find the current Hexagon ID:
		ueX := sc.GetUserByID(uint(userID)).X()
		ueY := sc.GetUserByID(uint(userID)).Y()
		currHex := hexMap.FindContainingHex(ueX, ueY)

		//2nd tier neighbours
		//2nd tier cells are in the array snIds
		sNeighs := hexMap.SecondNeighbours(currHex.ID)
		snIds := []uint{}
		for i := 0; i < len(sNeighs); i++ {
			snIds = append(snIds, sNeighs[i].ID)
		}

		//The 2nd tier cells are assigned a code; Each cell's code ia matched with the current cell's code to find if that is a FR3 cells and stored in the array frCells
		frCells := []uint{}
		code := [19]uint{2, 1, 3, 1, 3, 2, 1, 3, 2, 1, 3, 2, 1, 3, 2, 1, 2, 1, 3}
		currFreqCode := code[currHex.ID]
		for k := 0; k < len(code); k++ {
			for j := 0; j < len(snIds); j++ {
				if currFreqCode == code[snIds[j]] {
					frCells = append(frCells, snIds[j])
				}
			}
			break
		}
		// Append the current cell as well in all potential interfering cells
		frCells = append(frCells, currHex.ID)

		// The interfering stations array initialization
		bsIds = *new([]uint)
		for k := 0; k < len(frCells); k++ {
			bs := hexMap.FindContainedStations(frCells[k])
			for m := 0; m < len(bs); m++ {
				// Check if the BS's operator is enabled
				if p.OpEnableFlags[bs[m].OwnerOp().ID()] == true {
					bsIds = append(bsIds, bs[m].ID())
				}
			}
		}

	case "FFR":
		// FFR uses both FR1 and FR3; All UEs in the current node's cell are found.
		// Post SINR of all UEs in the current cell are found at mode FR1. Then the
		// top 'perc' % of users are assigned FR1. The remaining users are assigned FR3.
		perc := 50

		//Get x,y locations of UE and find the current Hexagon ID
		ueX := sc.GetUserByID(uint(userID)).X()
		ueY := sc.GetUserByID(uint(userID)).Y()
		currHex := hexMap.FindContainingHex(ueX, ueY)
		// Finding all UEs in the current cell
		rootUsers := hexMap.FindContainedUsers(currHex.ID)

		// Calculating the FR1 Post SINR for all UEs in the current cell, and store it in an array called postSinrs
		postSinrs := []float64{}
		for k := 0; k < len(rootUsers); k++ {
			values := SinrProfile(sc, hexMap, rootUsers[k].ID(), 0, &Params{FrMode: "FR1", Level: p.Level, OpEnableFlags: p.OpEnableFlags, IntCancellers: p.IntCancellers})
			postSinrs = append(postSinrs, values["post"].(float64))
		}

		//Sort the array values
		postSinrs, ind := sort(postSinrs)

		// Assign FR1 to UEs that lie in the top perc % of the power array
		found := false
		for j := 0; j < len(ind)*perc/100; j++ {
			if userID == rootUsers[ind[j]].ID() {
				found = true
				bsIds = intrStations(sc, hexMap, userID, &Params{FrMode: "FR1", Level: p.Level, OpEnableFlags: p.OpEnableFlags, IntCancellers: p.IntCancellers})
				break
			}
		}
		//Assign FR3 to the remaining UEs
		if found == false {
			bsIds = intrStations(sc, hexMap, userID, &Params{FrMode: "FR3", Level: p.Level, OpEnableFlags: p.OpEnableFlags, IntCancellers: p.IntCancellers})
		}

	case "AFFR":
		// AFFR uses both FR1 and FR3; All UEs from the current node's cell are found;
		// Post SINR of all UEs in the current cell are found at mode FR1. Then the
		// top 'th1' % of users are assigned FR1. The next 'th2' % of users are assigned FR3;
		// and the remaining are assigned FR3, but with only one operator.
		th1 := 50
		th2 := 40

		// Get x,y locations of UE and find the current cell's ID
		ueX := sc.GetUserByID(uint(userID)).X()
		ueY := sc.GetUserByID(uint(userID)).Y()
		currHex := hexMap.FindContainingHex(ueX, ueY)
		// Finding all the UEs in the current cell
		rootUsers := hexMap.FindContainedUsers(currHex.ID)

		// Calculating the FR1 Post SINR for all UEs in the current cell, and store it in an array called postSinrs
		postSinrs := []float64{}
		for k := 0; k < len(rootUsers); k++ {
			values := SinrProfile(sc, hexMap, rootUsers[k].ID(), 0, &Params{FrMode: "FR1", Level: p.Level, OpEnableFlags: p.OpEnableFlags, IntCancellers: p.IntCancellers})
			postSinrs = append(postSinrs, values["post"].(float64))
		}

		// Sort the power array
		postSinrs, ind := sort(postSinrs)

		x1 := len(ind) * th1 / 100
		x2 := len(ind) * th2 / 100
		// Assign FR1 to UEs that lie in the top 'th1' % of the power array
		found := false
		for j := 0; j < x1; j++ {
			if userID == rootUsers[ind[j]].ID() {
				found = true
				bsIds = intrStations(sc, hexMap, userID, &Params{FrMode: "FR1", Level: p.Level, OpEnableFlags: p.OpEnableFlags, IntCancellers: p.IntCancellers})
				break
			}
		}
		// Assign FR3 to UEs that lie in the next 'th2' % of the power array
		if found == false {
			for j := x1; j < (x1 + x2 - 1); j++ {
				if userID == rootUsers[ind[j]].ID() {
					found = false
					bsIds = intrStations(sc, hexMap, userID, &Params{FrMode: "FR3", Level: p.Level, OpEnableFlags: p.OpEnableFlags, IntCancellers: p.IntCancellers})
					break
				}
			}
		}
		// Assign FR3, but with only one operator from each FR3 cell to UEs that remain.
		if found == false {
			fmt.Printf("\nRare case reached!\n")
			var desOp = sc.GetUserByID(userID).CurrOp.ID()
			// for j := (x1 + x2 - 1); j < len(ind); j++ {
			// fmt.Printf("Loop reached, iteration: %v, usid: %v\n", ind[j], usid[ind[j]])
			// if userID == rootUsers[ind[j]].ID {
			// fmt.Println(" The selected UE ", userID, " follows FR3 with one operator")
			bsIdsAll := intrStations(sc, hexMap, userID, &Params{FrMode: "FR3", Level: p.Level, OpEnableFlags: p.OpEnableFlags, IntCancellers: p.IntCancellers})
			bsIds = *new([]uint)
			for i := 0; i < len(bsIdsAll); i++ {
				// fmt.Printf("Desired: %v, Iteration: %v", desOp, sc.GetStationByID(bsIdsAll[i]).OwnerOp())
				if sc.GetStationByID(bsIdsAll[i]).OwnerOp().ID() == desOp {
					bsIds = append(bsIds, bsIdsAll[i])
				}
			}
		}

	default:
		return nil
	}

	return bsIds
}
