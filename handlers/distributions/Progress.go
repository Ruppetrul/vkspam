package distributions

import "sync"

var (
	mutexes  = make(map[int]*sync.Mutex)
	percents = make(map[int]int)
)

func UpdateProgress(id int, percent int) {
	if _, exists := mutexes[id]; !exists {
		mutexes[id] = &sync.Mutex{}
	}

	mutexes[id].Lock()
	defer mutexes[id].Unlock()
	percents[id] = percent
}

func GetProgress(id int) int {
	mutexes[id].Lock()
	defer mutexes[id].Unlock()
	return percents[id]
}
