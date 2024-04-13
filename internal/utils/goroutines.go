package utils

import (
	"reflect"
	"sync"
)

func ActiveGoroutinesCount(wg *sync.WaitGroup) int {
	var mu sync.Mutex
	mu.Lock()
	count := wgCounter(wg)
	mu.Unlock()
	return count
}

func wgCounter(wg *sync.WaitGroup) int {
	counterField := reflect.ValueOf(wg).Elem().FieldByName("counter")
	return int(counterField.Int())
}
