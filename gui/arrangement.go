package gui

type Dimensions struct {
	X0 int
	X1 int
	Y0 int
	Y1 int
}

type Direction int

const (
	ROW Direction = iota
	COLUMN
)

type Container struct {
    Direction Direction

    Children []*Container

	Window string

	Size int
}
