package server

type Conn interface {
	Conn() interface{}
	IP() string
	Write(data []byte) error
	Read(data []byte) error
	Close() error
}
