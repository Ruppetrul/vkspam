package distributions

import (
	"sync"
	"vkspam/handlers/responses"
)

var (
	mutexes       = make(map[int]*sync.Mutex) // 0-100 - ok. -1 final -2 - error.
	percents      = make(map[int]int)
	statusMessage = make(map[int]string)
)

func UpdateProgress(id int, percent int, message string) {
	if _, exists := mutexes[id]; !exists {
		mutexes[id] = &sync.Mutex{}
	}

	mutexes[id].Lock()
	defer mutexes[id].Unlock()
	percents[id] = percent
	statusMessage[id] = message
}

func GetProgress(id int) *responses.ProcessDistributionResponse {
	if _, exists := mutexes[id]; !exists {
		return &responses.ProcessDistributionResponse{
			Progress: -1,
			Message:  "",
		}
	}
	mutexes[id].Lock()
	defer mutexes[id].Unlock()
	return &responses.ProcessDistributionResponse{
		Progress: percents[id],
		Message:  statusMessage[id],
	}
}

func DeleteProgress(id int) {
	if _, exists := mutexes[id]; !exists {
		mutexes[id] = &sync.Mutex{}
	}

	mutexes[id].Lock()
	defer mutexes[id].Unlock()
	delete(percents, id)
	delete(mutexes, id)
}
