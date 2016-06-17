package model

// Scenario encapsulates a physical scenario by maintaing a record of BaseStations, Users and Operators.
// It manages a loss table for each User with each BaseStation.
//
// Additionally, Scenario ensures error-free behaviour by rejecting duplicate ID's and
// provides a fast look-up of nodes by ID.
//
// Warning: The loss table look-ups are guranteed to work correctly provided ID's are consecutive and start from 0.
type Scenario struct {
	baseStations []*BaseStation
	users        []*User
	operators    []*Operator

	lossTable [][]float64

	baseStationMap map[uint]*BaseStation
	userMap        map[uint]*User
	operatorMap    map[uint]*Operator
}

// An O(1) fast lookup to retrieve BaseStations by ID
func (sc *Scenario) GetStationByID(id uint) *BaseStation {
	return sc.baseStationMap[id]
}

// An O(1) fast lookup to retrieve Users by ID
func (sc *Scenario) GetUserByID(id uint) *User {
	return sc.userMap[id]
}

// An O(1) fast lookup to retrieve Operators by ID
func (sc *Scenario) GetOperatorByID(id uint) *Operator {
	return sc.operatorMap[id]
}

// Gets a list of all BaseStations in the scenario. The array is free to modify as it is a light-weight copy of the actual.
func (sc *Scenario) BaseStations() []*BaseStation {
	bs := make([]*BaseStation, len(sc.baseStations))
	copy(bs, sc.baseStations)
	return bs
}

// Gets a list of all Users in the scenario. The array is free to modify as it is a light-weight copy of the actual.
func (sc *Scenario) Users() []*User {
	ue := make([]*User, len(sc.users))
	copy(ue, sc.users)
	return ue
}

// Gets a list of all Operators in the scenario. The array is free to modify as it is a light-weight copy of the actual.
func (sc *Scenario) Operators() []*Operator {
	op := make([]*Operator, len(sc.operators))
	copy(op, sc.operators)
	return op
}

// Retreives the loss values for a certain User. The index equals the ID of the BaseStation.
// The data returned is free to modify.
func (sc *Scenario) LossProfile(ueID uint) []float64 {
	losses := make([]float64, len(sc.baseStations))
	// Copying the data so that it cannot be modified from outside. (slices are reference types)
	copy(losses, sc.lossTable[ueID])
	return losses
}

// The loss in dBm between a User identified by ueID and a BaseStation identified by bsID.
func (sc *Scenario) Loss(ueID, bsID uint) float64 {
	return sc.lossTable[ueID][bsID]
}

// Move a User by displacement (dx, dy). This function also updates the corresponding loss values.
// WARNING: This method is incomplete. Do not use.
func (sc *Scenario) MoveUser(ueID uint, dx, dy float64) {
	ue := sc.GetUserByID(ueID)
	ue.x += dx
	ue.y += dy
	// Update loss data
	for i := 0; i < len(sc.baseStations); i++ {
		sc.lossTable[ueID][i] = float64(sc.lossTable[ueID][i]) // TODO: Calculate loss for new location
	}
}

// Private helper method: Add a new BaseStation safely by checking uniqueness.
func (sc *Scenario) addBaseStation(bs *BaseStation) bool {
	if _, exists := sc.baseStationMap[bs.ID()]; exists {
		return false
	}
	sc.baseStationMap[bs.ID()] = bs
	sc.baseStations = append(sc.baseStations, bs)
	return true
}

// Private helper method: Add a new User safely by checking uniqueness.
func (sc *Scenario) addUser(ue *User) bool {
	if _, exists := sc.userMap[ue.ID()]; exists {
		return false
	}
	sc.userMap[ue.ID()] = ue
	sc.users = append(sc.users, ue)
	return true
}

// Private helper method: Add a new Operator safely by checking uniqueness.
func (sc *Scenario) addOperator(op *Operator) bool {
	if _, exists := sc.operatorMap[op.ID()]; exists {
		return false
	}
	sc.operatorMap[op.ID()] = op
	sc.operators = append(sc.operators, op)
	return true
}

// Private: Constructor for a new Scenario. Use this only, and not the default initialization.
func newScenario() *Scenario {
	sc := new(Scenario)
	sc.baseStationMap = make(map[uint]*BaseStation)
	sc.userMap = make(map[uint]*User)
	sc.operatorMap = make(map[uint]*Operator)
	return sc
}
