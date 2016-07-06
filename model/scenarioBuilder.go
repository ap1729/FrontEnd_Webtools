package model

// The ScenarioBuilder type is a robust and safe factory to create new scenarios.
// It works by adding necessary node details and uses a "sealing" mechanism to prevent multiple references modifying data.
//
// The ScenarioBuilder is recommended and is the only way to make Scenario's currently.
//     1. Loss data can be imported only after all nodes have been added.
//     2. Unique ID's are enforced.
//     3. Copying a ScenarioBuilder object to another will not duplicate the scenario; however once a ScenarioBuilder is sealed, all its reference copies are sealed.
type ScenarioBuilder struct {
	scenario *Scenario
	lastBsId uint
	lastUeId uint
	isSealed *bool
}

// Checks if the ScenarioBuilder is sealed.
func (sb *ScenarioBuilder) IsSealed() bool {
	return *sb.isSealed
}

// The number of Users that have been added.
func (sb *ScenarioBuilder) NumUsersAdded() uint {
	return sb.lastUeId
}

// The number of BaseStations that have been added.
func (sb *ScenarioBuilder) NumStationsAdded() uint {
	return sb.lastBsId
}

// Adds an operator with the specified parameters. This function can only be used if the ScenarioBuilder is not sealed.
func (sb *ScenarioBuilder) AddOperator(opID uint) bool {
	if *sb.isSealed == true {
		return false
	}
	op := NewOperator(opID)
	err := sb.scenario.addOperator(op)
	if err == true {
		return false
	}
	return true
}

// Adds a node with the specified parameters. This function can only be used if the ScenarioBuilder is not sealed.
//
// NodeType specifies the type of node - "BS" and "UE" are supported.
func (sb *ScenarioBuilder) AddNode(nodeType string, x, y, ht float64, opID uint) bool {
	if *sb.isSealed == true {
		return false
	}
	op := sb.scenario.GetOperatorByID(opID)
	// If operator requested is not yet added to the scenario, fail
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
	// Unknown nodeType received, fail
	return false
}

// Check if the operator specified by id exists.
func (sb *ScenarioBuilder) OperatorExists(id uint) bool {
	return sb.scenario.GetOperatorByID(id) != nil
}

// Seal a ScenarioBuilder, to prevent further additions to nodes or operators or any changes.
// Once sealed, it cannot be undone, and all copies of the current object are also sealed internally.
//
// The lossOpt specifies if the loss table must be calculated ("calc") or imported ("import").
// In case of import, the second argument must contain all values as a (nUE x nBS) array.
//
// Note: General convention is to call Finalize() after Seal(), to retreive the Scenario object.
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
		sb.scenario.lossTable = make([][]float64, M)
		for i := 0; i < M; i++ {
			sb.scenario.lossTable[i] = make([]float64, N)
			for j := 0; j < N; j++ {
				sb.scenario.lossTable[i][j] = HataLoss(sb.scenario.BaseStations()[j].X(), sb.scenario.BaseStations()[j].Y(),
					sb.scenario.Users()[i].X(), sb.scenario.Users()[i].Y())
			}
		}
		goto successCase
	}

	return false

successCase:
	*sb.isSealed = true
	return true
}

// Constructor for ScenarioBuilder, use this only to instantiate.
func NewScenarioBuilder() *ScenarioBuilder {
	sealed := false
	return &ScenarioBuilder{scenario: newScenario(), lastBsId: 0, lastUeId: 0, isSealed: &sealed}
}

// Finalize returns the built Scenario object and destroys internal handles and resets to empty.
// Once Finalized, a ScenarioBuilder object is useless and can be deleted.
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
