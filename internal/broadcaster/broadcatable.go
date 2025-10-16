package broadcaster

import "errors"

var NotExists = errors.New("Not exists")

type UpdateFunction = func() error
type OnUpdated = func()

type ChangeHandlerInterface interface {
	OnChange(name string)
}

type Broadcastable interface {
	GetName() string
	Broadcast(replicaID int64, name string) error
	Update(data interface{}) error
	SetOnChanged(handler ChangeHandlerInterface)
}
