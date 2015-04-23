package chmap

//impl Imap
type Chmap struct {
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
	v := Chmap{make(map[string]interface{}), make(chan mission)}
	go v.init()
	return v
}

func (m *Chmap) Delete(k string) {
	m.request <- mission { DEL, kv{k, nil} }
}

func (m *Chmap) ErrorPut(k string, v interface{}) {
	m.kv[k] = v
}

func (m *Chmap) init() {
	for mission := range m.request {
		switch mission.mission_type {
		case DEL:
			delete(m.kv, mission.kv.k)
		case PUT:
			m.kv[mission.kv.k] = mission.kv.v
		}
	}
}

func (m *Chmap) Get(k string) interface{} {
	return m.kv[k]
}

func (m *Chmap) Put(k string, v interface{}) {
	m.request <- mission {PUT, kv{k, v}}
}

func (m *Chmap) Count() int {
	return len(m.kv)
}

func (m *Chmap) Contain(k string) bool {
	_, ok = m.kv[k]
	return ok
}
