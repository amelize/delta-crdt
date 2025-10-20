package object

type ObjectType uint8

const (
	ObjTypeCounter = iota
	ObjTypeText
	ObjTypeSet
	ObjTypeMap
)

type Replicatable interface {
	Broadcast() error
	Update(data interface{}) error
}

type Object struct {
	ObjectID   uint64
	Type       ObjectType
	InnterType Replicatable
}
