package service

import (
	"FrontEnd_WebTools/model"
	"encoding/csv"
	"os"
	"strconv"
)

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
