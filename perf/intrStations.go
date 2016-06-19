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
func intrStations(mode string, sc *model.Scenario, userID uint, params map[string]interface{}) []uint {
	var bsIds []uint

	switch mode {
	case "FR1":
		bsIds = make([]uint, len(sc.BaseStations()))
		for i := 0; i < len(sc.BaseStations()); i++ {
			bsIds[i] = uint(i)
		}
	case "FR3":

		hexMap := params["hexmap"].(*service.HexMap)
		//get x,y locations of UE
		uex := sc.GetUserByID(uint(userID)).X()
		uey := sc.GetUserByID(uint(userID)).Y()
		//fmt.Println("UE : ",userID, " x :",uex, " y : ",uey)
		//func to find the hex containing the UE
		currenthex := hexMap.FindContainingHex(uex, uey)
		//fmt.Println("currenthex: \n",currenthex)

		//cn is the current cell id
		cn := currenthex.ID

		//fmt.Printf("\n\nConent in map: %v\n\n", params)

		fmt.Println("UE : ", userID, " Hexagon id :", cn)

		/*
			//finding the strongest BS of the selected UE by calculating Lev 0 array, then use the co-ordinates of the BS to find the hexagon ID
				bsIds = make([]uint, len(sc.BaseStations()))
				for i := 0; i < len(sc.BaseStations()); i++ {
							bsIds[i] = uint(i)
						}
				losses, bsId := signalLossProfile(userID, sc, uint(0), bsIds)
				strongestBS:=bsId[0]
				fmt.Println("////in intrstns.go////\nLosses :",losses,"\n bsId : ",bsId)
				fmt.Println("strongestBS :",strongestBS)
				bsx:=sc.GetStationByID(uint(strongestBS)).X()
				bsy:=sc.GetStationByID(uint(strongestBS)).Y()
				fmt.Println("BSX:",bsx," BSY:",bsy)
				currenthex:=hexMap.FindContainingHex(bsx,bsy)
		*/

		//2nd tier neighbors
		//2nd tier cells are in the array snids
		sneigh := hexMap.SecondNeighbours(int(cn))
		snids := []uint{}
		for i := 0; i < len(sneigh); i++ {
			snids = append(snids, sneigh[i].ID)
		}
		//fmt.Println("2nd tier cells:",snids)

		//3 arrays containing the cell ids of same freqs.ie all cells in ar1 have the same freq
		ar1 := [7]uint{1, 3, 6, 9, 12, 15, 17}
		ar2 := [6]uint{0, 5, 8, 11, 14, 16}
		ar3 := [6]uint{2, 4, 7, 10, 13, 18}

		//array frc stores the fr3 cellids
		frc := []uint{}

		//to find if current id is in ar1,ar2 or ar3
		switch cn {
		case 1, 3, 6, 9, 12, 15, 17:
			for k := 0; k < len(ar1); k++ {
				for l := 0; l < len(snids); l++ {
					if snids[l] == ar1[k] {
						frc = append(frc, ar1[k])
					}
				}
			}
		case 0, 5, 8, 11, 14, 16:
			for k := 0; k < len(ar2); k++ {
				for l := 0; l < len(snids); l++ {
					if snids[l] == ar2[k] {
						frc = append(frc, ar2[k])
					}
				}
			}
		case 2, 4, 7, 10, 13, 18:
			for k := 0; k < len(ar3); k++ {
				for l := 0; l < len(snids); l++ {
					if snids[l] == ar3[k] {
						frc = append(frc, ar3[k])
					}
				}
			}

		}

		//fmt.Println("FR3 cells :",frc)
		//append the current cell also to the list of fr3 cells
		frc = append(frc, cn)
		//fmt.Println("FR3 cells :",frc)
		//array fr3bsno stores the BS IDS in fr3 cells
		//array fr3pow stores the SIR values of corresponding UEs in fr3 cells
		fr3bsno := []uint{}
		for k := 0; k < len(frc); k++ {
			//append each cell's bsid to fr3bsno
			for j := 0; j < len(sneigh); j++ {
				if sneigh[j].ID == frc[k] {
					//fmt.Println(&sneigh[j])
					bs1 := hexMap.FindContainedStations(&sneigh[j])
					for m := 0; m < len(bs1); m++ {
						id := bs1[m].ID()
						fr3bsno = append(fr3bsno, id)
					}
				}
			}
		}
		//appending the current cell's bs and pow
		bs2 := hexMap.FindContainedStations(currenthex)
		for m := 0; m < len(bs2); m++ {
			id := bs2[m].ID()
			fr3bsno = append(fr3bsno, id)
		}
		bsIds = make([]uint, len(fr3bsno))

		for p := 0; p < len(fr3bsno); p++ {
			bsIds[p] = fr3bsno[p]
		}
		//fmt.Println("Bsids :",bsIds)
		return bsIds

	case "FFR":
		fmt.Println("UE : ", userID)
		return nil

	case "AFFR":
	default:
		return nil
	}

	return bsIds
}
