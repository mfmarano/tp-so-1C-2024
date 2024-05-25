package io

import (
	"sync"
)

type IoMap struct {
	mutex sync.Mutex
	Ios map[string]IoConfig
}

type IoConfig struct {
	Ip string
	Port int
}

var IosMap *IoMap

func (ioMap *IoMap) AddConfig(name string, config IoConfig) {
	ioMap.mutex.Lock()
	ioMap.Ios[name] = config
	ioMap.mutex.Unlock()
}

func (ioMap *IoMap) GetConfig(name string) (IoConfig, bool) {
	ioMap.mutex.Lock()
	config, ok := ioMap.Ios[name]
	ioMap.mutex.Unlock()
	return config, ok
}