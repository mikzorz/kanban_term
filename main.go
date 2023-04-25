package main

func main() {

	s := newScreen()

	s.SetStyle(defStyle)

	quit := func() {
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	newNote("Note 1")
	newNote("Note 2")
	newNote("Note 3 with an extra long message just to test truncation")
	selected = 0

	updateLoop(s)
}
