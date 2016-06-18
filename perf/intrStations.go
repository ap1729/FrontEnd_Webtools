package perf

import "FrontEnd_WebTools/model"

// Returns a list of ID's that identify interfering BaseStations for a user in
// the scenario with given frequency-reuse mode.
//
// The params map is optional, and can be used to specify additional arguments
// that may be required for evaluation at a given frequency-reuse mode. If not
// needed, pass nil.
func intrStations(frMode string, sc *model.Scenario, userID uint, params map[string]interface{}) []uint {
	var bsIds []uint

	switch frMode {
	case "FR1":
		bsIds = make([]uint, len(sc.BaseStations()))
		for i := 0; i < len(sc.BaseStations()); i++ {
			bsIds[i] = uint(i)
		}
	case "FR3":
	case "FFR":
	case "AFFR":
	default:
		return nil
	}

	return bsIds
}
