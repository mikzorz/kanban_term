package main

type Box struct {
	x, y, w, h int
}

type BoundaryBox struct {
	x1, y1, x2, y2 int
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
