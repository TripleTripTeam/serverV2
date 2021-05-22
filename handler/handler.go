package handler

type Handler struct {
	gbChan chan Coord
	Coords map[int]map[int]map[int]Coord
}

func NewHandler(gbChan chan Coord) *Handler {
	return &Handler{
		gbChan: gbChan,
		Coords: make(map[int]map[int]map[int]Coord),
	}
}
