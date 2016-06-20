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
		hexMap := params["hexmap"].(*service.HexMap)
		intrCancelCount := (params["intcnc"].(uint))
		//get x,y locations of UE
		uex := sc.GetUserByID(uint(userID)).X()
		uey := sc.GetUserByID(uint(userID)).Y()

		//func to find the hex containing the UE
		currenthex := hexMap.FindContainingHex(uex, uey)

		//cn is the current cell id
		cn := currenthex.ID
		fmt.Println("UE : ", userID, " Hexagon id :", cn)

		//finding the UEs in a cell
		usid := []uint{}
		us := hexMap.FindContainedUsers(currenthex)
		for j := 0; j < len(us); j++ {
			id := us[j].ID()
			usid = append(usid, id)
		}

		//fmt.Println("Users :",usid)
		//fmt.Println("len :",len(usid))

		//finding FR1 Post SINR for all UEs in the current cell, and store it in an array called posarr
		posarr := []float64{}
		for k := 0; k < len(usid); k++ {
			// bsIds1 := intrStations("FR1", sc, uint(k), nil)
			// losses, bsId := signalLossProfile(uint(k), sc, 1, bsIds1)
			//fmt.Println("losses: ",losses)
			//fmt.Println("bsid:",bsId)
			// op := make([]uint, len(bsId))
			// for i := 0; i < len(bsId); i++ {
			// 	losses[i] += 46.0
			// 	op[i] = sc.GetStationByID(bsId[i]).OwnerOp().ID()
			// }
			// arr := []float64{}
			// arr = sinr(losses, intrCancelCount) //The func sinr() returns an array of three values..arr[0] gives Pre SINR, arr[1] gives Post SINR and arr[2] gives the ROI values.
			//Since only the Post SINR is needed, arr[1] is considered.
			// posarr = append(posarr, arr[1])
			values := SinrProfile(sc, "FR1", usid[k], 1, intrCancelCount, 1, nil)
			posarr = append(posarr, values["SINR"].([]float64)[1])
		}
		//fmt.Println("POST SINR :",posarr)
		//fmt.Println("len(sinr): ",len(posarr))

		posarr, ind := sort(posarr)

		/*l1:=filter(posarr,usid)
		los,in:=sort(l1)
		usid1:=make([]uint,len(in))
		var x uint
		for j:=0;j<len(in);j++{
			x=usid[in[j]]
			usid1[j]=sc.GetUserByID(x).ID()
		}*/
		//fmt.Println("--After sort :--")
		//fmt.Println("PoSarr: ",posarr)
		//fmt.Println(usid)
		//
		// var perc float64
		// var percval float64
		// //give the values as %..i.e.if 60%, give as 60
		perc := 50
		// //fmt.Println("Length of ue1 :",lenue1)
		// percval = (perc / 100.0) * float64(len(usid))
		// top1 := int(math.Floor(percval))
		// rem1 := len(usid) - int(top1)
		//fmt.Println("Perc:",perc," percval :",percval)
		//fmt.Println(" Top1:",top1," Rem1 :",rem1)
		// fmt.Println("Out of ", len(usid), " UEs in the cell ", cn, " ,")
		// fmt.Println("the no.of cells that follow FR1 :", top1)
		// fmt.Println("the no.of cells that follow FR3 :", rem1)
		// fr1ues := []uint{}
		// fr3ues := []uint{}
		// for i := 0; i < len(usid); i++ {
		// 	if i < top1 {
		// 		fr1ues = append(fr1ues, usid[ind[i]])
		// 	} else {
		// 		fr3ues = append(fr3ues, usid[ind[i]])
		// 	}
		// }
		//fmt.Println("FR1 : ",fr1ues,"len :",len(fr1ues))
		//fmt.Println("FR3 : ",fr3ues,"len :",len(fr3ues))
		t := 0
		for j := 0; j < len(ind)*perc/100; j++ {
			if userID == usid[ind[j]] {
				t = 1
				fmt.Println(" The selected UE ", userID, " follows FR1")
				bsIds = intrStations("FR1", sc, userID, params)
			}
		}
		if t == 0 {
			fmt.Println(" The selected UE ", userID, " follows FR3")
			bsIds = intrStations("FR3", sc, userID, params)
			break
		}

	case "AFFR":
	default:
		return nil
	}

	return bsIds
}
