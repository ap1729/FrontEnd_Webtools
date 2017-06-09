package service

import (
	"FrontEnd_WebTools/model"
	"encoding/csv"
	"os"
	"strconv"
	"fmt"
)

// Reads a CSV file containing locations of nodes, and populates the ScenarioBuilder.
// The CSV file is expected to have a format as follows:
//
//     type, x, y
//
// where type is a string who's first two letters specify the type of node ("BS", "UE")
// and remaining characters specify the operator ID. (Ex. "BS1", "UE3")
func ReadNodes(sb *model.ScenarioBuilder, locFilePath string) bool {
	file, err := os.Open(locFilePath)
	if err != nil {
		fmt.Println(err)
		return false
	}

	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		return false
	}
	// Temporary fix
	firstLine := true
	for _, row := range records {

		// Temporary fix
		if firstLine == true {
			firstLine = false
			continue
		}
		nodeType := row[0][0:2]
		//fmt.Println(nodeType)

		        op, err := strconv.ParseInt(string(row[3]), 10, 64)
				if err != nil {
					fmt.Println(err)
					return false
				}
                 //if not assigned any operator ,op 10 is error operator
				if(op==-1){
				 op=10
				}
				opId := uint(op)
        lvlbs0,err := strconv.ParseInt(string(row[4]),10,64)
        if err != nil {
        	fmt.Println(err)
			return false
		}
		lvlbs1,err := strconv.ParseInt(string(row[5]),10,64)
        if err != nil {
        	fmt.Println(err)
			return false
		}

		if !sb.OperatorExists(opId) {
			sb.AddOperator(opId)
		}
		x, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			fmt.Println(err)
			return false
		}
		y, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			fmt.Println(err)
			return false
		}
		err2 := sb.AddNode(nodeType, x, y, 0, opId,lvlbs0,lvlbs1)
		if err2 == false {
			return false
		}
	}

	return true
}

// Takes an unsealed ScenarioBuilder and imports loss data from a CSV file into the builder.
// This function seals the ScenarioBuilder on completion.
//
// The CSV file is to have M rows and N columns, where M is the number of Users and N is the number of BaseStations.
// Each row is an array of N values separated by comma.
func ReadLossTable(sb *model.ScenarioBuilder, lossFilePath string) bool {
	if sb.IsSealed() == true {
		return false
	}
	file, err := os.Open(lossFilePath)
	if err != nil {
		return false
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return false
	}
	lossTable := make([][]float64, sb.NumUsersAdded())
	for i, _ := range records {
		lossTable[i] = make([]float64, sb.NumStationsAdded())
		for j, val := range records[i] {
			loss, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return false
			}
			lossTable[i][j] = loss
		}
	}
	sb.Seal("import", lossTable)
	return true
}
