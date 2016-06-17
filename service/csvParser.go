package service

import (
	"FrontEnd_WebTools/model"
	"encoding/csv"
	"os"
	"strconv"
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
		return false
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return false
	}
	for _, row := range records {
		nodeType := row[0][0:2]
		opId, err := strconv.ParseUint(string(row[0][2]), 10, 64)
		if err != nil {
			return false
		}
		if !sb.OperatorExists(uint(opId)) {
			sb.AddOperator(uint(opId))
		}
		x, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			return false
		}
		y, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return false
		}
		err2 := sb.AddNode(nodeType, x, y, 0, uint(opId))
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
