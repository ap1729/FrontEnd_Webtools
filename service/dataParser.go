package service

import (
	"FrontEnd_WebTools/model"
	"encoding/csv"
	"os"
	"strconv"
)

// Generate a scenario object with all nodes and operators
func ReadNodes(filePath string) (*model.Scenario, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var bsId, ueId uint = 0, 0
	var sc = model.NewScenario()
	for _, row := range records {
		nodeType := row[0][0:2]
		opId, err := strconv.ParseUint(string(row[0][2]), 10, 64)
		if err != nil {
			return nil, err
		}
		op := model.NewOperator(uint(opId))
		sc.AddOperator(op)
		x, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			return nil, err
		}
		y, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, err
		}
		if nodeType == "BS" {
			sc.AddBaseStation(model.NewBaseStation(bsId, x, y, 0, op))
			bsId++
		} else if nodeType == "UE" {
			sc.AddUser(model.NewUser(ueId, x, y, 0, op))
			ueId++
		} else {
			return sc, nil
		}
	}
	return sc, nil
}

// Load the loss values into the LossTable of the scenario reference
func ImportLossTable(filePath string, destSc *model.Scenario) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}
	destSc.LossTable = make([][]float64, len(destSc.Users))
	for i, _ := range records {
		destSc.LossTable[i] = make([]float64, len(destSc.BaseStations))
		for j, val := range records[i] {
			loss, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return err
			}
			destSc.LossTable[i][j] = loss
		}
	}
	return nil
}
