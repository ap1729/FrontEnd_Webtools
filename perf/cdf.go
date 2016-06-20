package perf

import (
	"FrontEnd_WebTools/model"
	"fmt"
	"math"
	s "sort"
)

//ramanan cdf changes begin global variables and structures declaration

var numUeCdf int = 100 // total number of UE for CDF calculation is hardcoded
var minOfMin, maxOfMax float64

type cdfL0L1Ret struct { //structure to hold SINR CDF
	prsDbXArr  []float64 //pre procesing SINR in db x-axis
	prsCdfYArr []float64 //pre cdf y-axis
	posDbXArr  []float64 //post procesing SINR in db x-axis
	posCdfYArr []float64 //post cdf y-axis
}
type calCdfRet struct { //structure to hold SINR CDF
	sinrDbXArr []float64 //SINR in db x-axis
	cdfYArr    []float64 //cdf y-axis
}

/*
type prefxPosfx struct{			// pre and post fix 0s and 1s to CDF
      cdfValues []float64
      num0Prefix float64
      num1Postfix float64
}
*/

func CDF(sc *model.Scenario, frMode string, intrCancelCount uint) map[string]interface{} {
	var cdfL0Obj, cdfL1Obj cdfL0L1Ret
	returnData := map[string]interface{}{}
	fmt.Println("CDF func reached")

	cdfL0Obj = cdfL0L1(sc, uint(0), intrCancelCount, frMode)
	cdfL1Obj = cdfL0L1(sc, uint(1), intrCancelCount, frMode)

	var sinrMinCollectArr = []float64{} //collect all minimum in X
	sinrMinCollectArr = append(sinrMinCollectArr, cdfL0Obj.prsDbXArr[0])
	sinrMinCollectArr = append(sinrMinCollectArr, cdfL0Obj.posDbXArr[0])
	sinrMinCollectArr = append(sinrMinCollectArr, cdfL1Obj.prsDbXArr[0])
	sinrMinCollectArr = append(sinrMinCollectArr, cdfL1Obj.posDbXArr[0])
	s.Float64s(sinrMinCollectArr) //function which sorts data in ascending order
	minOfMin = sinrMinCollectArr[0]

	var sinrMaxCollectArr = []float64{} //collect all maximum in X
	sinrMaxCollectArr = append(sinrMaxCollectArr, cdfL0Obj.prsDbXArr[0]+float64(len(cdfL0Obj.prsDbXArr)-1))
	sinrMaxCollectArr = append(sinrMaxCollectArr, cdfL0Obj.posDbXArr[0]+float64(len(cdfL0Obj.posDbXArr)-1))
	sinrMaxCollectArr = append(sinrMaxCollectArr, cdfL1Obj.prsDbXArr[0]+float64(len(cdfL1Obj.prsDbXArr)-1))
	sinrMaxCollectArr = append(sinrMaxCollectArr, cdfL1Obj.posDbXArr[0]+float64(len(cdfL1Obj.posDbXArr)-1))
	s.Float64s(sinrMaxCollectArr) //function which sorts data in ascending order
	maxOfMax = sinrMaxCollectArr[3]

	//combining x array
	var combineSinrX = []float64{}
	tempX := minOfMin
	for int(tempX) <= int(maxOfMax) {
		combineSinrX = append(combineSinrX, tempX)
		tempX = tempX + float64(1)
	}

	returnData["sinr"] = combineSinrX
	returnData["pre0"] = pre0Post1(cdfL0Obj.prsCdfYArr, cdfL0Obj.prsDbXArr[0], float64(len(cdfL0Obj.prsDbXArr)-1))
	returnData["post0"] = pre0Post1(cdfL0Obj.posCdfYArr, cdfL0Obj.posDbXArr[0], float64(len(cdfL0Obj.posDbXArr)-1))
	returnData["pre1"] = pre0Post1(cdfL1Obj.prsCdfYArr, cdfL1Obj.prsDbXArr[0], float64(len(cdfL1Obj.prsDbXArr)-1))
	returnData["post1"] = pre0Post1(cdfL1Obj.posCdfYArr, cdfL1Obj.posDbXArr[0], float64(len(cdfL1Obj.posDbXArr)-1))
	return returnData
}

func pre0Post1(cdfArr []float64, firstEleXArr float64, lenMinus1 float64) []float64 {
	zeroSlice := []float64{0}
	num0Prefix := firstEleXArr - minOfMin
	num1Postfix := maxOfMax - (firstEleXArr + lenMinus1)
	for i := 0; i < int(num0Prefix); i++ {
		cdfArr = append(zeroSlice, cdfArr...)
	}
	for j := 0; j < int(num1Postfix); j++ {
		cdfArr = append(cdfArr, float64(1))
	}
	return cdfArr
}

func cdfL0L1(temp *model.Scenario, levelValue uint, intfCancelCount uint, frMode string) cdfL0L1Ret {
	var cdfL0L1RetObj cdfL0L1Ret
	var calCdfObj calCdfRet
	var prsPosRoiArr = []float64{} //array containing pre , post processing SINR and ROI values
	var prsArrDb = []float64{}     //array variable to store pre processing SINR for number of UEs considered
	var posArrDb = []float64{}     //array variable to store post processing SINR for number of UEs considered
	for userId := 0; userId < numUeCdf; userId++ {
		intStatIds := intrStations(frMode, temp, uint(userId), nil)
		losses, bsId := signalLossProfile(uint(userId), temp, levelValue, intStatIds)

		for i := 0; i < len(bsId); i++ {
			losses[i] += 46
		}
		prsPosRoiArr = sinr(losses, intfCancelCount)
		prsArrDb = append(prsArrDb, prsPosRoiArr[0])
		posArrDb = append(posArrDb, prsPosRoiArr[1])
	}

	//>>>>>>>>>>>>>>>>>> for pre processing SINR <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<//
	calCdfObj = calCdf(prsArrDb)
	cdfL0L1RetObj.prsDbXArr = append(cdfL0L1RetObj.prsDbXArr, calCdfObj.sinrDbXArr...)
	cdfL0L1RetObj.prsCdfYArr = append(cdfL0L1RetObj.prsCdfYArr, calCdfObj.cdfYArr...)
	//>>>>>>>>>>>>>>>>>>>>>> for post processing SINR <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<//
	calCdfObj = calCdf(posArrDb)
	cdfL0L1RetObj.posDbXArr = append(cdfL0L1RetObj.posDbXArr, calCdfObj.sinrDbXArr...)
	cdfL0L1RetObj.posCdfYArr = append(cdfL0L1RetObj.posCdfYArr, calCdfObj.cdfYArr...)
	return cdfL0L1RetObj
}

func calCdf(sinrValuesArr []float64) calCdfRet {
	var calCdfRetObj calCdfRet
	//fmt.Println("\n pre_sinr_arr_dB_before_sorting \n",sinr_values_arr,"\n")
	s.Float64s(sinrValuesArr) //function which sorts data in ascending order
	var minSinrDB int
	var maxSinrDB int
	minSinrDB = int(math.Floor(sinrValuesArr[0]))
	maxSinrDB = int(math.Ceil(sinrValuesArr[numUeCdf-1]))

	sinrThreshold := minSinrDB       // %%%%%%%%%%%%%%threshold is kept as integer as of now%%%%%%%%%%%%%%%%
	for sinrThreshold <= maxSinrDB { // x axis range
		calCdfRetObj.sinrDbXArr = append(calCdfRetObj.sinrDbXArr, float64(sinrThreshold))
		sinrCountUe := 0
		for j := 0; j < numUeCdf; j++ { //  calculateCDF
			if sinrValuesArr[j] <= float64(sinrThreshold) {
				sinrCountUe = sinrCountUe + 1
			}
		}
		//sinr_cdf = float64(sinrCountUe)/float64(numUeCdf)
		//fmt.Println("\n check values \n",pre_sinr_cdf)
		calCdfRetObj.cdfYArr = append(calCdfRetObj.cdfYArr, float64(sinrCountUe)/float64(numUeCdf))
		sinrThreshold = sinrThreshold + 1
	}
	return calCdfRetObj
}
