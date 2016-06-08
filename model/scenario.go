package model

type Scenario struct {
	// The exported variables are not read-only, awaiting language feature updates (hopefully!)
	BaseStations []BaseStation
	Users        []User
	Operators    []Operator
	// The loss table is valid under the assumption that all arrays are ordered by index
	LossTable [][]float64

	baseStationMap map[uint]BaseStation
	userMap        map[uint]User
	operatorMap    map[uint]Operator
}

// Provide a fast lookup to retrieve objects by ID
func (sc *Scenario) GetStationByID(id uint) BaseStation {
	return sc.baseStationMap[id]
}
func (sc *Scenario) GetUserByID(id uint) User {
	return sc.userMap[id]
}
func (sc *Scenario) GetOperatorByID(id uint) Operator {
	return sc.operatorMap[id]
}

// Methods to add new nodes safely by checking uniqueness
func (sc *Scenario) AddBaseStation(bs *BaseStation) bool {
	if _, exists := sc.baseStationMap[bs.ID()]; exists {
		return false
	}
	sc.baseStationMap[bs.ID()] = *bs
	sc.BaseStations = append(sc.BaseStations, *bs)
	return true
}
func (sc *Scenario) AddUser(ue *User) bool {
	if _, exists := sc.userMap[ue.ID()]; exists {
		return false
	}
	sc.userMap[ue.ID()] = *ue
	sc.Users = append(sc.Users, *ue)
	return true
}
func (sc *Scenario) AddOperator(op *Operator) bool {
	if _, exists := sc.operatorMap[op.ID()]; exists {
		return false
	}
	sc.operatorMap[op.ID()] = *op
	sc.Operators = append(sc.Operators, *op)
	return true
}

// "Constructor" for a new Scenario
func NewScenario() *Scenario {
	sc := new(Scenario)
	sc.baseStationMap = make(map[uint]BaseStation)
	sc.userMap = make(map[uint]User)
	sc.operatorMap = make(map[uint]Operator)
	return sc
}
