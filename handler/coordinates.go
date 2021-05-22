package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
)

const (
	StartAngle = -0.33529248
	Step       = 0.016297472
)

type inputCoord struct {
	X    float64   `json:"x"`
	Y    float64   `json:"y"`
	Z    float64   `json:"z"`
	Phi  float64   `json:"phi"`
	Dist []float64 `json:"dist"`
}

type Coord struct {
	X float64
	Y float64
	Z float64
}

func (h *Handler) SendCoords(w http.ResponseWriter, r *http.Request) {

	var input inputCoord
	body, err := ioutil.ReadAll(r.Body)
	if err != nil || len(body) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println(input)

	for j, dist := range input.Dist {
		var coord Coord
		teta := math.Pi/2 - (StartAngle + Step*float64(j))
		coord.X = dist*math.Sin(teta)*math.Cos(input.Phi*0.017) + input.X
		coord.Y = dist*math.Sin(teta)*math.Sin(input.Phi*0.017) + input.Y
		coord.Y = dist*math.Cos(teta) + input.Z
		h.gbChan <- coord
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (h *Handler) CoordProcess() {
	for {
		coord := <-h.gbChan
		h.CoordFilter(coord)

		fmt.Println(h.Coords)
	}
}

func (h *Handler) CoordFilter(coord Coord) {
	xSearch := int(coord.X * 100)
	ySearch := int(coord.Y * 100)
	zSearch := int(coord.Z * 100)
	if _,ok := h.Coords[xSearch][ySearch][zSearch]; !ok{
		zMap := map[int]Coord{
			zSearch: coord,
		}
	
		yMap := map[int]map[int]Coord{
			ySearch: zMap,
		}
		
		h.Coords[xSearch] = yMap 
	}

	
}
