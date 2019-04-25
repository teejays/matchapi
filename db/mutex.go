package db

import (
	"sync"
)

var lockMap = map[string]*sync.RWMutex{
	UserCollection: &sync.RWMutex{},
	LikeCollection: &sync.RWMutex{},
}

func lock(collection string) {
	lockMap[collection].Lock()
}

func unlock(collection string) {
	lockMap[collection].Unlock()
}

func rlock(collection string) {
	lockMap[collection].RLock()
}

func runlock(collection string) {
	lockMap[collection].RUnlock()
}
