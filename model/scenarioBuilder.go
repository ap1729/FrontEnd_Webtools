package model

type ScenarioBuilder struct {
	scenario *Scenario
	lastBsId uint
	lastUeId uint

	opExists map[uint]bool
	isSealed *bool
}

func (sb *ScenarioBuilder) IsSealed() bool {
	return *sb.isSealed
}
func (sb *ScenarioBuilder) NumUsersAdded() uint {
	return sb.lastUeId
}
func (sb *ScenarioBuilder) NumStationsAdded() uint {
	return sb.lastBsId
}

func (sb *ScenarioBuilder) AddOperator(opID uint) bool {
	if *sb.isSealed == true {
		return false
	}
	op := NewOperator(opID)
	err := sb.scenario.addOperator(op)
	if err == true {
		return false
	}
	sb.opExists[opID] = true
	return true
}
func (sb *ScenarioBuilder) AddNode(nodeType string, x, y, ht float64, opID uint) bool {
	if *sb.isSealed == true {
		return false
	}
	op := sb.scenario.GetOperatorByID(opID)
	// If operator requested is not yet added to the scenario
	if op == nil {
		return false
	}
	if nodeType == "BS" {
		bs := NewBaseStation(sb.lastBsId, x, y, ht, op)
		sb.lastBsId++
		return sb.scenario.addBaseStation(bs)
	}
	if nodeType == "UE" {
		ue := NewUser(sb.lastUeId, x, y, ht, op)
		sb.lastUeId++
		return sb.scenario.addUser(ue)
	}
	// Unknown nodeType received
	return false
}
func (sb *ScenarioBuilder) OperatorExists(id uint) bool {
	return sb.opExists[id]
}

func (sb *ScenarioBuilder) Seal(lossOpt string, lossTable [][]float64) bool {
	if *sb.isSealed == true {
		return false
	}
	M := len(sb.scenario.users)
	N := len(sb.scenario.baseStations)
	if lossOpt == "import" {
		if lossTable != nil && len(lossTable) != M {
			return false
		}
		sb.scenario.lossTable = make([][]float64, M)
		for i := 0; i < M; i++ {
			sb.scenario.lossTable[i] = make([]float64, N)
			if lossTable[i] != nil && len(lossTable[i]) != N {
				return false
			}
			for j := 0; j < N; j++ {
				sb.scenario.lossTable[i][j] = lossTable[i][j]
			}
		}
		goto successCase
	}
	if lossOpt == "calc" {
		for i := 0; i < M; i++ {
			sb.scenario.lossTable[i] = make([]float64, N)
			for j := 0; j < N; j++ {
				sb.scenario.lossTable[i][j] = 0 // Code to calculate losses
			}
		}
		goto successCase
	}

	return false

successCase:
	*sb.isSealed = true
	return true
}

func NewScenarioBuilder() *ScenarioBuilder {
	sealed := false
	return &ScenarioBuilder{scenario: newScenario(), lastBsId: 0, lastUeId: 0, isSealed: &sealed, opExists: make(map[uint]bool)}
}

// Finalize returns the built Scenario object and destroys internal handles and resets to empty
func (sb *ScenarioBuilder) Finalize() *Scenario {
	if *sb.isSealed != true {
		return nil
	}
	sc := sb.scenario
	sb.scenario = nil
	sb.lastBsId = 0
	sb.lastUeId = 0
	return sc
}
