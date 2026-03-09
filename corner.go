package main

type Corner int

const (
	TOP_LEFT Corner = iota
	TOP_RIGHT
	BOTTOM_RIGHT
	BOTTOM_LEFT
)

var cornerName = map[Corner]string{
	TOP_LEFT:     "TOP_LEFT",
	TOP_RIGHT:    "TOP_RIGHT",
	BOTTOM_RIGHT: "BOTTOM_RIGHT",
	BOTTOM_LEFT:  "BOTTOM_LEFT",
}

func (c Corner) String() string {
	if cornerString, ok := cornerName[c]; ok {
		return cornerString
	} else {
		return "???"
	}
}
