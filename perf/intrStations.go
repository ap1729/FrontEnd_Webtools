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
		//2nd tier neighbors
		//2nd tier cells are in the array snids
		sneigh := hexMap.SecondNeighbours(cn)
		snids := []uint{}
		for i := 0; i < len(sneigh); i++ {
			snids = append(snids, sneigh[i].ID)
		}
		//fmt.Println("2nd tier cells:",snids)

		//array frc stores the fr3 cellids
		frc := []uint{}
		code:=[19] uint{2,1,3,1,3,2,1,3,2,1,3,2,1,3,2,1,2,1,3}
		f:=code[cn]
		//fmt.Println("f: ",f)
		for k:=0;k<len(code);k++{
			for j:=0;j<len(snids);j++{
				if f==code[snids[j]] {
					frc=append(frc,snids[j])
					}
			}
			break
		}
		//fmt.Println("FR3 cells :",frc)
		//append the current cell also to the list of fr3 cells
		frc = append(frc, cn)
		fmt.Println("FR3 cells :",frc)
		//array fr3bsno stores the BS IDS in fr3 cells
		fr3bsno := []uint{}
		for k := 0; k < len(frc); k++ {
			//append each cell's bsid to fr3bsno
			bs1 := hexMap.FindContainedStations(frc[k])
					for m := 0; m < len(bs1); m++ {
						id := bs1[m].ID()
						fr3bsno = append(fr3bsno, id)
					}
				}

		bsIds = make([]uint, len(fr3bsno))
		for p := 0; p < len(fr3bsno); p++ {
			bsIds[p] = fr3bsno[p]
		}
		fmt.Println("Bsids :",bsIds, "\n len:",len(bsIds) )
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
