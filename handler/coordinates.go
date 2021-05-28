package handler

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"math"
	"net/http"
	"os"
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

	fmt.Println("sendCoords")
	err = json.Unmarshal(body, &input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println(input.Dist)
	for j, dist := range input.Dist {
		var coord Coord
		teta := math.Pi/2 - (StartAngle + Step*float64(j))
		coord.X = dist * math.Sin(teta) * math.Cos(input.Phi)
		coord.Y = dist * math.Sin(teta) * math.Sin(input.Phi)
		coord.Z = dist * math.Cos(teta)
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
	}
}

func (h *Handler) CoordFilter(coord Coord) {
	xSearch := int(coord.X * 100)
	ySearch := int(coord.Y * 100)
	zSearch := int(coord.Z * 100)
	if _, ok := h.Coords[xSearch][ySearch][zSearch]; !ok {
		zMap := map[int]Coord{
			zSearch: coord,
		}

		yMap := map[int]map[int]Coord{
			ySearch: zMap,
		}

		h.Coords[xSearch] = yMap
	}

}
func (h *Handler) Print(w http.ResponseWriter, r *http.Request) {
	csvfile, err := os.Create("output.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer csvfile.Close()

	writer := csv.NewWriter(csvfile)
	writer.Write([]string{"x1", "y1", "z1"})
	writer.Write([]string{"-1", "-1", "-1"})
	writer.Write([]string{"1", "1", "1"})
	writer.Write([]string{"0", "0", "0"})

	for _, val1 := range h.Coords {
		for _, val2 := range val1 {
			for _, coord := range val2 {
				str := []string{fmt.Sprintf("%f", coord.X), fmt.Sprintf("%f", coord.Y), fmt.Sprintf("%f", coord.Z)}
				writer.Write(str)
			}
		}
	}

	writer.Flush()

	path := "../public/index.html"

	//создаем html-шаблон
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func (h *Handler) Output(w http.ResponseWriter, r *http.Request) {
	path2 := "../cmd/output.csv"
	http.ServeFile(w, r, path2)
}
