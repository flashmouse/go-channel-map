package chmap

import "hash/fnv"

//impl Imap
type Chmap struct {
	//	kv map[string]interface{}
	//	request chan mission
	innerMaps []innerMap
	shardNum  int
}

type innerMap struct {
	kv map[string]interface{}
	request chan mission
}

var (
	PUT int = 1
	DEL int = 2
)

type kv struct {
	k string
	v interface{}
}

type mission struct{
	mission_type int
	kv           kv
}

func NewMap() Chmap {
	//	v := Chmap{make(map[string]interface{}), make(chan mission)}
	v := Chmap{}
	v.shardNum = 16
	v.innerMaps = make([]innerMap, v.shardNum, v.shardNum)
	for i := 0; i < v.shardNum; i++ {
		v.innerMaps[i] = innerMap{make(map[string]interface{}), make(chan mission)}
		go v.innerMaps[i].init()
	}
	return v
}

//func NewMap(int shardNum) Chmap{
//	return nil
//}

func (m *Chmap) getShard(k string) uint {
	h := fnv.New32()
	h.Write([]byte(k))
	return uint(h.Sum32()) % uint(m.shardNum)
}

func (m *Chmap) Delete(k string) {
	m.innerMaps[m.getShard(k)].request <- mission { DEL, kv{k, nil} }
}

func (m *innerMap) init() {
	for mission := range m.request {
		switch mission.mission_type {
		case DEL:
			delete(m.kv, mission.kv.k)
		case PUT:
			m.kv[mission.kv.k] = mission.kv.v
		}
	}
}

func (m *Chmap) Get(k string) (interface{}, bool) {
	v , ok := m.innerMaps[m.getShard(k)].kv[k]
	return v, ok
}

func (m *Chmap) Put(k string, v interface{}) {
	m.innerMaps[m.getShard(k)].request <- mission {PUT, kv{k, v}}
}

func (m *Chmap) Count() int {
	re := 0
	for i := 0; i < m.shardNum; i++ {
		re += len(m.innerMaps[i].kv)
	}
	return re
}

func (m *Chmap) Contain(k string) bool {
	_, ok := m.innerMaps[m.getShard(k)].kv[k]
	return ok
}
