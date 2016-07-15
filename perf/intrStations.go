package perf

import (
	"FrontEnd_WebTools/model"
	"FrontEnd_WebTools/service"
	"errors"
	"fmt"
)

// TODO: Major optimization opportunity exists in intrStations; for FFR and AFFR
// where every user in the cell needs to be compared, a cache mechanism could be
// implemented to save the resultant calculation frMode in a hash table.

// Returns a list of ID's that identify interfering BaseStations for a user in
// the scenario with given frequency-reuse mode.
func intrStations(sc *model.Scenario, hexMap *service.HexMap, userID uint, p *Params) ([]uint, error) {

	// Handling argument nil exception
	if sc == nil || hexMap == nil || p == nil {
		return nil, errors.New(ARG_NIL)
	}

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
		// FR3 sends a list of BS in FR3 cells that match the operator(s) enabled.
		// Get x,y locations of UE and find the current Hexagon ID:
		ueX := sc.GetUserByID(uint(userID)).X()
		ueY := sc.GetUserByID(uint(userID)).Y()
		currHex := hexMap.FindContainingHex(ueX, ueY)

		if currHex == nil {
			return nil, fmt.Errorf("The userID %d requested is not associated with any cell.", userID)
		}

		// 2nd tier neighbours
		// 2nd tier cells are in the array snIds
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
		// top 'th' % of users are assigned FR1. The remaining users are assigned FR3.
		th := 50

		//Get x,y locations of UE and find the current Hexagon ID
		ueX := sc.GetUserByID(uint(userID)).X()
		ueY := sc.GetUserByID(uint(userID)).Y()
		currHex := hexMap.FindContainingHex(ueX, ueY)

		if currHex == nil {
			return nil, fmt.Errorf("The userID %d requested is not associated with any cell.", userID)
		}

		// Finding all UEs in the current cell
		rootUsers := hexMap.FindContainedUsers(currHex.ID)

		// Calculating the FR1 Post SINR for all UEs in the current cell, and store it in an array called postSinrs
		postSinrs := []float64{}
		for k := 0; k < len(rootUsers); k++ {
			values, err := SinrProfile(sc, hexMap, rootUsers[k].ID(), 0, &Params{FrMode: "FR1", Level: p.Level, OpEnableFlags: p.OpEnableFlags, IntCancellers: p.IntCancellers})
			if err != nil {
				return nil, fmt.Errorf("Failed to evaluate user profiles, at userID: %d", rootUsers[k].ID())
			}
			postSinrs = append(postSinrs, values["post"].(float64))
		}

		// Sort the array values
		postSinrs, ind := sort(postSinrs)

		var err error
		// Assign FR1 to UEs that lie in the top 'th' % of the power array
		found := false
		for j := 0; j < len(ind)*th/100; j++ {
			if userID == rootUsers[ind[j]].ID() {
				found = true
				bsIds, err = intrStations(sc, hexMap, userID, &Params{FrMode: "FR1", Level: p.Level, OpEnableFlags: p.OpEnableFlags, IntCancellers: p.IntCancellers})
				break
			}
		}
		// Assign FR3 to the remaining UEs
		if found == false {
			bsIds, err = intrStations(sc, hexMap, userID, &Params{FrMode: "FR3", Level: p.Level, OpEnableFlags: p.OpEnableFlags, IntCancellers: p.IntCancellers})
		}

		if err != nil {
			return nil, err
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

		if currHex == nil {
			return nil, fmt.Errorf("The userID %d requested is not associated with any cell.", userID)
		}

		// Finding all the UEs in the current cell
		rootUsers := hexMap.FindContainedUsers(currHex.ID)

		// Calculating the FR1 Post SINR for all UEs in the current cell, and store it in an array called postSinrs
		postSinrs := []float64{}
		for k := 0; k < len(rootUsers); k++ {
			values, err := SinrProfile(sc, hexMap, rootUsers[k].ID(), 0, &Params{FrMode: "FR1", Level: p.Level, OpEnableFlags: p.OpEnableFlags, IntCancellers: p.IntCancellers})
			if err != nil {
				return nil, fmt.Errorf("Failed to evaluate user profiles, at userID: %d", rootUsers[k].ID())
			}
			postSinrs = append(postSinrs, values["post"].(float64))
		}

		// Sort the power array
		postSinrs, ind := sort(postSinrs)

		x1 := len(ind) * th1 / 100
		x2 := len(ind) * th2 / 100
		var err error
		// Assign FR1 to UEs that lie in the top 'th1' % of the power array
		found := false
		for j := 0; j < x1; j++ {
			if userID == rootUsers[ind[j]].ID() {
				found = true
				bsIds, err = intrStations(sc, hexMap, userID, &Params{FrMode: "FR1", Level: p.Level, OpEnableFlags: p.OpEnableFlags, IntCancellers: p.IntCancellers})
				break
			}
		}
		// Assign FR3 to UEs that lie in the next 'th2' % of the power array
		if found == false {
			for j := x1; j < (x1 + x2 - 1); j++ {
				if userID == rootUsers[ind[j]].ID() {
					found = false
					bsIds, err = intrStations(sc, hexMap, userID, &Params{FrMode: "FR3", Level: p.Level, OpEnableFlags: p.OpEnableFlags, IntCancellers: p.IntCancellers})
					break
				}
			}
		}
		// Assign FR3, but with only one operator from each FR3 cell to UEs that remain.
		if found == false {
			fmt.Printf("\nRare case reached!\n")
			var desOp = sc.GetUserByID(userID).CurrOp.ID()
			var bsIdsAll []uint
			bsIdsAll, err = intrStations(sc, hexMap, userID, &Params{FrMode: "FR3", Level: p.Level, OpEnableFlags: p.OpEnableFlags, IntCancellers: p.IntCancellers})
			if err == nil {
				bsIds = *new([]uint)
				for i := 0; i < len(bsIdsAll); i++ {
					// fmt.Printf("Desired: %v, Iteration: %v", desOp, sc.GetStationByID(bsIdsAll[i]).OwnerOp())
					if sc.GetStationByID(bsIdsAll[i]).OwnerOp().ID() == desOp {
						bsIds = append(bsIds, bsIdsAll[i])
					}
				}
			}
		}

		if err != nil {
			return nil, err
		}

	default:
		return nil, errors.New("Unsupported FR mode requested in parameters.")
	}

	return bsIds, nil
}
