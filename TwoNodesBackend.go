package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

var records map[int]Node // Entries of all the nodes
var BSNodes []int        // stores the index of BSnodes

// type TestRecord struct {
// 	col1 string
// 	node string
// 	x    float64
// 	y    float64
// 	m1   float64
// 	m2   float64
// }

// var allRecords []TestRecord //global
// var BSRecords []TestRecord

//object for storing values

/// E.g for sending from browser
// Obj={id:"32",Column1:30,Column2:-45,Column3:-34,Column4:-40}
// $.post("http://localhost:8080/update",JSON.stringify(Obj),"json")
// type rowdata struct {
// 	Node, Col                          int
// 	Column1, Column2, Column3, Column4 float64
// 	X1, Y1, X2, Y2                     float64

// 	//node is indication of Basestation or node
// 	//Column1 to 4 are x,y,m1 and m2 values originally
// 	//X1,X2,Y1,Y2 are x and y coordinates of basestation locations
// }

func handlerroute(w http.ResponseWriter, r *http.Request) {

	// w.Header().Add("Access-Control-Allow-Origin", "*")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Only POST method is supported"))
		return
	}
	log.Println("Request Method is ", r.Method)
	// log.Println("Request is originated from  ", r.RemoteAddr)
	// log.Println("Request is originated URL  ", r.RequestURI)
	// log.Println("Request Headers", r.Header)

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("body is ", string(data))

		var node Node
		jerr := json.Unmarshal(data, &node)

		if jerr != nil {
			log.Println("Unmarshalling error ", jerr)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		udpateData := EvaluateMore(node)
		txbytes, _ := json.Marshal(udpateData)

		nbytes, werr := w.Write(txbytes)
		_ = nbytes
		if werr != nil {
			log.Println("I got some error while writing back", werr)
		} else {
			log.Println("Sent this  ", string(txbytes))
			// log.Printf("Successfully returned %d bytes", nbytes)
		}

	}

}

func distanceBetween(srcid, destid int) float64 {
	s := records[srcid]
	d := records[destid]
	return (s.x - d.x) // not the correct expression
}

func EvaluateMore(node Node) []Node {
	// log.Printf("I received this to process %#v", r)
	// Actual algo goes here ....
	var result []Node
	log.Println("Received a BS or NODES ??", node)

	if node.nodetype != 0 {

		// var a, b, dist1, dist2, BSX1, BSY1, BSX2, BSY2 float64
		// a = (r.Column1) / 30
		// b = (r.Column2 + 5) / 30
		// BSX1 = (r.X1) / 30
		// BSY1 = (r.Y1 + 5) / 30
		// BSX2 = (r.X2) / 30
		// BSY2 = (r.Y2 + 5) / 30

		// dist1 = ((a-BSX1)*(a-BSX1) + (b-BSY1)*(b-BSY1))
		// dist2 = ((a-BSX2)*(a-BSX2) + (b-BSY2)*(b-BSY2))

		// r.Column4 = (0.2) * ((1 / dist1) + (1 / dist2))
		// r.Column3 = 500 + 80*math.Log(math.Abs(dist1*dist2))

		// if r.Column3 < 0 {
		// 	r.Column3 = 0
		// }

		// allRecords[r.Col-len(BSRecords)].x = r.X1
		// allRecords[r.Col-len(BSRecords)].y = r.Y1
		// allRecords[r.Col-len(BSRecords)].m1 = r.Column3
		// allRecords[r.Col-len(BSRecords)].m2 = r.Column4

		for bsid := range BSNodes {
			dist1 := distanceBetween(bsid, node.ID)
			// find the metrics m1 based on the distance
			log.Println("Distance is  ", dist1, "from  ", bsid)
			// update the value of m1, m2 etc.. of the node
			// node1.m1 = 0 // some function of distance

		}
		// update back in the record
		records[node.ID] = node
		result = append(result, node) /// send an array of result of length 1

	} else {
		//Basesation stuff of saving session and BS coords
		//log.Println("BASESTATION!!!")
		bsid := node.ID
		for key, val := range result {

			// process all nodes except the current as well as other BS nodes
			if val.nodetype != 0 {
				log.Println("Processing Node ID ", val.ID)
				dist1 := distanceBetween(bsid, val.ID)
				// find the metrics m1 based on the distance
				log.Println("Distance is  ", dist1, "from  ", bsid)

				/// update this node's property
				// val.m1=0
				// val.m2=0
				result[key] = val

				/// add to the result
				result = append(result, val)

			}
		}

		// update the value of m1, m2 etc.. of the node
		// node1.m1 = 0 // some function of distance

		//Need to return bunch of data
	}
	return result

}

type Node struct {
	ID       int
	nodetype int
	x        float64
	y        float64
	m1       float64
	m2       float64
}

func (n *Node) Parse(str []string) bool {
	var ferr, err error
	n.ID, err = strconv.Atoi(str[0])
	if err != nil {
		ferr = err
	}
	n.nodetype, err = strconv.Atoi(str[1])
	if err != nil {
		ferr = err
	}
	n.x, err = strconv.ParseFloat(str[2], 64)
	if err != nil {
		ferr = err
	}
	n.y, err = strconv.ParseFloat(str[3], 64)
	if err != nil {
		ferr = err
	}
	n.m1, err = strconv.ParseFloat(str[4], 64)
	if err != nil {
		ferr = err
	}
	n.m2, err = strconv.ParseFloat(str[5], 64)
	if err != nil {
		ferr = err
	}
	if ferr != nil {
		return false
	} else {
		return true
	}
}

func LoadRecords(fname string) {
	csvfile, err := os.Open(fname)
	records = make(map[int]Node)
	BSNodes = []int{} // clearing the entries
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer csvfile.Close()
	reader := csv.NewReader(csvfile)
	headers, err := reader.Read() /// assuming first row is header
	fmt.Println("HEADERS ", headers)
	var cols []string
	var cnt int
	for err == nil {
		//record, err := reader.Read() /// assuming first row is header
		cols, err = reader.Read()
		if err == nil {
			// fmt.Println(cnt, record)
			var node Node
			if node.Parse(cols) {
				// Successfull Node parsed
				records[node.ID] = node
				cnt++
				if node.nodetype == 0 {
					BSNodes = append(BSNodes, node.ID)
				}
			}

		}

	}
	// log.Println(records)
	log.Println(BSNodes)

}

func main() {
	LoadRecords("testfile5.csv")

	log.Println("Started Server at :8080")
	http.HandleFunc("/update", handlerroute)
	http.Handle("/", http.FileServer(http.Dir(".")))
	err1 := http.ListenAndServe(":8080", nil)
	if err1 != nil {
		log.Println("Error ", err1)
	}
}
