package model

type Scenario struct {
	baseStations []*BaseStation
	users        []*User
	operators    []*Operator

	// The loss table is valid under the assumption that ID's are consecutive and start from 0.
	lossTable [][]float64

	baseStationMap map[uint]*BaseStation
	userMap        map[uint]*User
	operatorMap    map[uint]*Operator
}

// Provide a fast lookup to retrieve objects by ID
func (sc *Scenario) GetStationByID(id uint) *BaseStation {
	return sc.baseStationMap[id]
}
func (sc *Scenario) GetUserByID(id uint) *User {
	return sc.userMap[id]
}
func (sc *Scenario) GetOperatorByID(id uint) *Operator {
	return sc.operatorMap[id]
}

// Getter Methods for all arrays - This implementation ensures immutability
// (at the cost of a little performance)
func (sc *Scenario) BaseStations() []*BaseStation {
	bs := make([]*BaseStation, len(sc.baseStations))
	copy(bs, sc.baseStations)
	return bs
}
func (sc *Scenario) Users() []*User {
	ue := make([]*User, len(sc.users))
	copy(ue, sc.users)
	return ue
}
func (sc *Scenario) Operators() []*Operator {
	op := make([]*Operator, len(sc.operators))
	copy(op, sc.operators)
	return op
}

// Access to loss data in an immutable way. The data returned is free to modify.
func (sc *Scenario) LossProfile(ueID uint) []float64 {
	losses := make([]float64, len(sc.baseStations))
	copy(losses, sc.lossTable[ueID])
	return losses
}
func (sc *Scenario) Loss(ueID, bsID uint) float64 {
	return sc.lossTable[ueID][bsID]
}

// Moving a user equipment, and changing corresponding loss values
func (sc *Scenario) MoveUser(ueID uint, dx, dy float64) {
	ue := sc.GetUserByID(ueID)
	ue.x += dx
	ue.y += dy
	// Update loss data
	for i := 0; i < len(sc.baseStations); i++ {
		sc.lossTable[ueID][i] = float64(sc.lossTable[ueID][i]) // Calculate loss for new location
	}
}

// Private helper methods to add new nodes safely by checking uniqueness
func (sc *Scenario) addBaseStation(bs *BaseStation) bool {
	if _, exists := sc.baseStationMap[bs.ID()]; exists {
		return false
	}
	sc.baseStationMap[bs.ID()] = bs
	sc.baseStations = append(sc.baseStations, bs)
	return true
}
func (sc *Scenario) addUser(ue *User) bool {
	if _, exists := sc.userMap[ue.ID()]; exists {
		return false
	}
	sc.userMap[ue.ID()] = ue
	sc.users = append(sc.users, ue)
	return true
}
func (sc *Scenario) addOperator(op *Operator) bool {
	if _, exists := sc.operatorMap[op.ID()]; exists {
		return false
	}
	sc.operatorMap[op.ID()] = op
	sc.operators = append(sc.operators, op)
	return true
}

// Private constructor for a new Scenario
func newScenario() *Scenario {
	sc := new(Scenario)
	sc.baseStationMap = make(map[uint]*BaseStation)
	sc.userMap = make(map[uint]*User)
	sc.operatorMap = make(map[uint]*Operator)
	return sc
}

// Provide API called "MoveUE" etc.
// Provide API to access loss data
