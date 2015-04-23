package chmap



type Imap interface {
	Get(k string) interface {}
	Put(k string, v interface{})
	Delete(k string)
	Count() int
	Contain(k string) bool
}
