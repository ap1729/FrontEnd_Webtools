package service

import (
	"math"
	"test/model"
)

type HexMap struct {
	hexagons       []Hexagon
	adjMat         [][]int
	memberStations map[Hexagon]([]*model.BaseStation)
}

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

func (hm *HexMap) AssociateStations(stations []model.BaseStation) int {
	missCount := 0
	flag := false
	for i := 0; i < len(stations); i++ {
		flag = false
		for j := 0; j < len(hm.hexagons); j++ {
			if hm.hexagons[j].Contains(stations[i].X(), stations[i].Y()) == true {
				hm.memberStations[hm.hexagons[j]] = append(hm.memberStations[hm.hexagons[j]], &stations[i])
				flag = true
			}
		}
		if flag == false {
			missCount++
		}
	}
	return missCount
}

func (hm *HexMap) FirstNeighbours(root uint) []Hexagon {
	var neighs []Hexagon
	for i := 0; i < len(hm.hexagons); i++ {
		if hm.adjMat[root][i] == 1 {
			neighs = append(neighs, hm.hexagons[i])
		}
	}
	return neighs
}

func (hm *HexMap) SecondNeighbours(root int) []Hexagon {
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
	for i := 0; i < N; i++ {
		if hm.adjMat[root][i] == 0 && sum[i] > 0 && root != i {
			neighs = append(neighs, hm.hexagons[i])
		}
	}
	return neighs
}

func (hm *HexMap) FindContainingHex(x, y float64) *Hexagon {
	for i := 0; i < len(hm.hexagons); i++ {
		if hm.hexagons[i].Contains(x, y) == true {
			return &(hm.hexagons[i])
		}
	}
	return nil
}

func (hm *HexMap) FindContainedStations(hex *Hexagon) []*model.BaseStation {
	return hm.memberStations[*hex]
}

// Helper function for internal usage
func euclideanDist(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1))
}
