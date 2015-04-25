package chmap

type Imap interface {
	Get(k string) (interface{},bool)
	Put(k string, v interface{})
	Delete(k string)
	Count() int
	Contain(k string) bool
}
