package perf

import (
	"FrontEnd_WebTools/model"
	"FrontEnd_WebTools/service"
	"fmt"
)

// Returns a list of ID's that identify interfering BaseStations for a user in
// the scenario with given frequency-reuse mode.
//
// The params map is optional, and can be used to specify additional arguments
// that may be required for evaluation at a given frequency-reuse mode. If not
// needed, pass nil.
func intrStations(mode string, sc *model.Scenario, userID uint, opEnable []bool, level uint, params map[string]interface{}) []uint {
	var bsIds []uint

	switch mode {
	case "FR1":
		//FR1 sends the list of BS in cells that match the operator(s) received
		bsIds = *new([]uint)
		for i := uint(0); i < uint(len(sc.BaseStations())); i++ {
			if opEnable[sc.GetStationByID(i).OwnerOp().ID()] == true {
				bsIds = append(bsIds, i)
			}
		}

	case "FR3":
		//FR3 sends a list of BS in FR3 cells that match the operator(s) received
		hexMap := params["hexmap"].(*service.HexMap)

		//Get x,y locations of UE and find the current Hexagon ID
		uex := sc.GetUserByID(uint(userID)).X()
		uey := sc.GetUserByID(uint(userID)).Y()
		currenthex := hexMap.FindContainingHex(uex, uey)

		//cn is the current cell id
		cn := currenthex.ID
		// fmt.Println("UE : ", userID, " Hexagon id :", cn)

		//2nd tier neighbors
		//2nd tier cells are in the array snids
		sneigh := hexMap.SecondNeighbours(cn)
		snids := []uint{}
		for i := 0; i < len(sneigh); i++ {
			snids = append(snids, sneigh[i].ID)
		}

		//The 2nd tier cells are assigned a code; Each cell's code ia matched with the current cell's code to find if that is a FR3 cells and stored in the array frc
		frc := []uint{}
		code := [19]uint{2, 1, 3, 1, 3, 2, 1, 3, 2, 1, 3, 2, 1, 3, 2, 1, 2, 1, 3}
		f := code[cn]
		for k := 0; k < len(code); k++ {
			for j := 0; j < len(snids); j++ {
				if f == code[snids[j]] {
					frc = append(frc, snids[j])
				}
			}
			break
		}
		frc = append(frc, cn)
		// fmt.Println("FR3 cells :", frc)

		//array bsIds stores the BS IDS of the required operators in fr3 cells
		bsIds = *new([]uint)
		for k := 0; k < len(frc); k++ {
			bs := hexMap.FindContainedStations(frc[k])
			for m := 0; m < len(bs); m++ {
				//check if the BS is of the required operator
				if opEnable[bs[m].OwnerOp().ID()] == true {
					bsIds = append(bsIds, bs[m].ID())
				}
			}
		}

	case "FFR":
		//FFR uses both FR1 and FR3; All UEs from the current node's cell are found; Post SINR of all UEs in the current cell are found using FR1. Then the top 50% of users are assigned FR1. The remaining 50% of users are assigned FR3.

		//params is modified to contain Hexagon details,the no.of Interference cancellers and level.
		hexMap := params["hexmap"].(*service.HexMap)
		intrCancelCount := (params["intcnc"].(uint))

		//Get x,y locations of UE and find the current Hexagon ID
		uex := sc.GetUserByID(uint(userID)).X()
		uey := sc.GetUserByID(uint(userID)).Y()
		currenthex := hexMap.FindContainingHex(uex, uey)

		//cn is the current cell id
		cn := currenthex.ID
		fmt.Println("UE : ", userID, " Hexagon id :", cn)

		//Finding all UEs in the current cell
		usid := []uint{}
		us := hexMap.FindContainedUsers(cn)
		for j := 0; j < len(us); j++ {
			id := us[j].ID()
			usid = append(usid, id)
		}

		//Calculating the FR1 Post SINR for all UEs in the current cell, and store it in an array called posarr
		posarr := []float64{}
		for k := 0; k < len(usid); k++ {
			values := SinrProfile(sc, "FR1", usid[k], level, intrCancelCount, 1, opEnable, nil)
			posarr = append(posarr, values["post"].(float64))
		}
		//Sort the array values
		posarr, ind := sort(posarr)

		//Fix a threshold value
		perc := 50

		//Assign FR1 to UEs that lie in the top perc % of the power array
		t := 0
		for j := 0; j < len(ind)*perc/100; j++ {
			if userID == usid[ind[j]] {
				t = 1
				fmt.Println(" The selected UE ", userID, " follows FR1")
				bsIds = intrStations("FR1", sc, userID, opEnable, level, params)
			}
		}
		//Assign FR3 to the remaining UEs
		if t == 0 {
			fmt.Println(" The selected UE ", userID, " follows FR3")
			bsIds = intrStations("FR3", sc, userID, opEnable, level, params)
			break
		}

	case "AFFR":
		//AFFR uses both FR1 and FR3; All UEs from the current node's cell are found; Post SINR of all UEs in the current cell are found using FR1. Then the top 50% of users are assigned FR1. The next 40% of users are assigned FR3; the last 10% are assigned FR3, but with only one operator.

		//params is modified to contain Hexagon details,the no.of Interference cancellers and level.
		hexMap := params["hexmap"].(*service.HexMap)
		intrCancelCount := (params["intcnc"].(uint))

		//Get x,y locations of UE and find the current cell's ID
		uex := sc.GetUserByID(uint(userID)).X()
		uey := sc.GetUserByID(uint(userID)).Y()
		currenthex := hexMap.FindContainingHex(uex, uey)

		//cn is the current cell id
		cn := currenthex.ID
		fmt.Println("UE : ", userID, " Hexagon id :", cn)

		//Finding all the UEs in the current cell
		usid := []uint{}
		us := hexMap.FindContainedUsers(cn)
		for j := 0; j < len(us); j++ {
			id := us[j].ID()
			usid = append(usid, id)
		}

		//Calculating the FR1 Post SINR for all UEs in the current cell, and store it in an array called posarr
		posarr := []float64{}
		for k := 0; k < len(usid); k++ {
			values := SinrProfile(sc, "FR1", usid[k], level, intrCancelCount, 1, opEnable, nil)
			posarr = append(posarr, values["post"].(float64))
		}

		//Sort the power array
		posarr, ind := sort(posarr)

		// Fix two threholds
		th1 := 50
		th2 := 40
		x1 := len(ind) * th1 / 100
		x2 := len(ind) * th2 / 100

		//Assign FR1 to UEs that lie in the top th1 % of the power array
		t := 0
		for j := 0; j < x1; j++ {
			if userID == usid[ind[j]] {
				t = 1
				fmt.Println(" The selected UE ", userID, " follows FR1")
				bsIds = intrStations("FR1", sc, userID, opEnable, level, params)
				break
			}
		}

		//Assign FR3 to UEs that lie in the next th2 % of the power array
		if t == 0 {
			for j := (x1); j < (x1 + x2 - 1); j++ {
				if userID == usid[ind[j]] {
					t = 1
					fmt.Println(" The selected UE ", userID, " follows FR3")
					bsIds = intrStations("FR3", sc, userID, opEnable, level, params)
					break
				}
			}
		}

		//Assign FR3, but with only one operator from each FR3 cell to UEs that  remain.
		if t == 0 {
			fmt.Printf("\nRare case reached!\n")
			fmt.Printf("User ID: %v, x1: %v, x2: %v", userID, x1, x2)
			var desop uint
			desop = sc.GetUserByID(userID).CurrOp.ID()
			fmt.Printf("Dest Op: %v, Max users: %v, Max ind: %v", desop, len(usid), len(ind))

			fmt.Printf("User ID's: %v", usid)

			for j := (x1 + x2 - 1); j < len(ind); j++ {
				fmt.Printf("Loop reached, iteration: %v, usid: %v\n", ind[j], usid[ind[j]])
				if userID == usid[ind[j]] {
					fmt.Println(" The selected UE ", userID, " follows FR3 with one operator")
					bsIdsAll := intrStations("FR3", sc, userID, opEnable, level, params)
					bsIds = *new([]uint)
					for i := 0; i < len(bsIdsAll); i++ {
						fmt.Printf("Desired: %v, Iteration: %v", desop, sc.GetStationByID(bsIdsAll[i]).OwnerOp())
						if sc.GetStationByID(bsIdsAll[i]).OwnerOp().ID() == desop {
							bsIds = append(bsIds, bsIdsAll[i])
						}
					}
					break
				}
			}
		}
	default:
		return nil
	}

	return bsIds
}
