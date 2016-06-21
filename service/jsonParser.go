package service

import "FrontEnd_WebTools/model"

func PackageScenario(sc *model.Scenario) map[string]interface{} {
	data := map[string]interface{}{}

	bs := make([]map[string]interface{}, len(sc.BaseStations()))
	for i := 0; i < len(sc.BaseStations()); i++ {
		bs[i] = map[string]interface{}{}
		bs[i]["id"] = sc.BaseStations()[i].ID()
		bs[i]["x"] = sc.BaseStations()[i].X()
		bs[i]["y"] = sc.BaseStations()[i].Y()
		bs[i]["opid"] = sc.BaseStations()[i].OwnerOp().ID()
	}

	ue := make([]map[string]interface{}, len(sc.Users()))
	for i := 0; i < len(sc.Users()); i++ {
		ue[i] = map[string]interface{}{}
		ue[i]["id"] = sc.Users()[i].ID()
		ue[i]["x"] = sc.Users()[i].X()
		ue[i]["y"] = sc.Users()[i].Y()
		ue[i]["opid"] = sc.Users()[i].DefaultOp().ID()
	}

	op := make([]map[string]interface{}, len(sc.Operators()))
	for i := 0; i < len(sc.Operators()); i++ {
		op[i] = map[string]interface{}{}
		op[i]["id"] = sc.Operators()[i].ID()
	}

	data["basestats"] = bs
	data["users"] = ue
	data["opers"] = op
	return data

}
