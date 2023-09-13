package cash

import "sync"

var Cash sync.Map

func GetCash() *sync.Map {
	return &Cash
}
