package distributions

import "sync"

var (
	mutexes  = make(map[int]*sync.Mutex) // 0-100 - ok. -1 - not found. -2 - error.
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
	if _, exists := mutexes[id]; !exists {
		return -1
	}
	mutexes[id].Lock()
	defer mutexes[id].Unlock()
	return percents[id]
}

func DeleteProgress(id int) {
	if _, exists := mutexes[id]; !exists {
		mutexes[id] = &sync.Mutex{}
	}

	mutexes[id].Lock()
	defer mutexes[id].Unlock()
	percents[id] = -1
}
