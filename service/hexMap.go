package service

import (
	"FrontEnd_WebTools/model"
	"math"
)

// TODO: This struct requires improvements in its API, especially FindContainedStations().
// HexMap models a cell map, that contains a map of Hexagonal cells.
type HexMap struct {
	hexagons       []Hexagon
	adjMat         [][]int
	memberStations map[Hexagon]([]*model.BaseStation)
}

// Constructor to generate a cell map, with suitable ID's assigned to each hexagon.
// This is the recommended way to generate multi-tier maps as found in practice.
//
// Params: 'sideLength' is the length of the side of the regular hexagonal cell,
// while the 'radius' is a natural number specifying the number of tiers to generate.
func NewHexMap(sideLength float64, radius uint) *HexMap {
	N := 2*radius - 1
	normLength := sideLength * math.Sqrt(3) / 2
	var cx, cy float64
	// Coordinates of the top-left hexagon
	var cxi float64 = -1*normLength*float64(N-1) + (float64((radius+1)%2) * normLength)
	var cyi float64 = float64(radius-1) * math.Sqrt(3) * normLength
	// Select those hexagons that lie radially around the center up to the extent radius
	var index uint = 0
	var hexagons []Hexagon = *new([]Hexagon)
	memberMap := map[Hexagon]([]*model.BaseStation){} // Initializing map
	for i := uint(0); i < N; i++ {
		for j := uint(0); j < N; j++ {
			cx = cxi + float64(j)*normLength*2 - float64(i%2)*normLength
			cy = cyi - float64(i)*math.Sqrt(3)*normLength
			if euclideanDist(cx, cy, 0, 0) <= float64(radius-1)*2*normLength+0.001 {
				hexagons = append(hexagons, Hexagon{CenterX: cx, CenterY: cy, Side: sideLength, ID: index})
				memberMap[hexagons[index]] = []*model.BaseStation{} // Initializing map
				index++
			}
		}
	}
	// Define adjacency matrix
	adjMat := make([][]int, index)
	for i := uint(0); i < index; i++ {
		adjMat[i] = make([]int, index)
		for j := uint(0); j < index; j++ {
			if i == j {
				adjMat[i][j] = 0
				continue
			}
			if euclideanDist(hexagons[i].CenterX, hexagons[i].CenterY, hexagons[j].CenterX, hexagons[j].CenterY) <= 2*normLength+0.001 {
				adjMat[i][j] = 1
			}
		}
	}
	return &HexMap{hexagons: hexagons, adjMat: adjMat, memberStations: memberMap}
}

// Parses a list of BaseStations and creates internal references to which BaseStation lies in which Hexagon.
//
// BUG: Each call appends the existing knowledge of BaseStation locations with the passed array of stations.
// If care is not exercised, this may lead to duplicate references.
func (hm *HexMap) AssociateStations(stations []*model.BaseStation) int {
	missCount := 0
	flag := false
	for i := 0; i < len(stations); i++ {
		flag = false
		for j := 0; j < len(hm.hexagons); j++ {
			if hm.hexagons[j].Contains(stations[i].X(), stations[i].Y()) == true {
				hm.memberStations[hm.hexagons[j]] = append(hm.memberStations[hm.hexagons[j]], stations[i])
				flag = true
			}
		}
		if flag == false {
			missCount++
		}
	}
	return missCount
}

// Returns a list of hexagons that are first neighbours to the root hexagon specified by its ID.
func (hm *HexMap) FirstNeighbours(root uint) []Hexagon {
	var neighs []Hexagon
	for i := 0; i < len(hm.hexagons); i++ {
		if hm.adjMat[root][i] == 1 {
			neighs = append(neighs, hm.hexagons[i])
		}
	}
	return neighs
}

// Returns a list of second-tier neighbours to the root hexagon specified by its ID.
//
// Simply put, these are the neighbours of neighbours of the root hexagon.
func (hm *HexMap) SecondNeighbours(root uint) []Hexagon {
	sum := make([]int, len(hm.adjMat))
	N := len(hm.hexagons)
	for i := 0; i < N; i++ {
		if hm.adjMat[root][i] == 1 {
			// Loop over adjacencies of neighbour i of root
			for j := 0; j < N; j++ {
				// Accumulate all adjacencies to subtly absorb duplicates
				sum[j] = sum[j] + hm.adjMat[i][j]
			}
		}
	}
	var neighs []Hexagon
	// A non-zero element in sum is a second tier neighbour if it is not connected to root
	for i := uint(0); i < uint(N); i++ {
		if hm.adjMat[root][i] == 0 && sum[i] > 0 && root != i {
			neighs = append(neighs, hm.hexagons[i])
		}
	}
	return neighs
}

// Finds the hexagon in the map that contains the point (x, y). If not found, the function returns nil.
func (hm *HexMap) FindContainingHex(x, y float64) *Hexagon {
	for i := 0; i < len(hm.hexagons); i++ {
		if hm.hexagons[i].Contains(x, y) == true {
			return &(hm.hexagons[i])
		}
	}
	return nil
}

// Retreive all BaseStations that are contained in the specified hexagon.
//
// This function essentially retreives the stations as parsed by AssociateStations().
func (hm *HexMap) FindContainedStations(hex *Hexagon) []*model.BaseStation {
	return hm.memberStations[*hex]
}

// Private helper function: Euclidean distance between two points (x1, y1) and (x2, y2).
func euclideanDist(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1))
}
