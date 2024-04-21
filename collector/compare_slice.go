package collector

import "github.com/buarki/find-castles/castle"

func slicesWithSameContent(x, y []castle.Model) bool {
	if len(x) != len(y) {
		return false
	}
	diff := make(map[castle.Model]int, len(x))
	for _, _x := range x {
		diff[_x]++
	}
	for _, _y := range y {
		if _, ok := diff[_y]; !ok {
			return false
		}
		diff[_y]--
		if diff[_y] == 0 {
			delete(diff, _y)
		}
	}
	return len(diff) == 0
}
