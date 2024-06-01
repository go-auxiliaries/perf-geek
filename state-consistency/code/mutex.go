package code

import "sync"

type StateMutex struct {
	state State
	mu    sync.RWMutex
}

func (d *StateMutex) Init() *StateMutex {
	return &StateMutex{}
}

func (d *StateMutex) AddFruit(fruit string) {
	d.mu.Lock()
	d.state.fruits = append(d.state.fruits, fruit)
	if fruit == "apple" {
		d.state.numberOfApples++
	}
	d.mu.Unlock()
}

func (d *StateMutex) GetState() State {
	d.mu.RLock()
	val := d.state
	d.mu.RUnlock()
	return val
}
