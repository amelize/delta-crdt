package broadcaster

type UpdatedHandlerInterface interface {
	OnUpdate()
}

type ChangeHandlerInterface interface {
	OnChange(name string)
}

type Broadcastable interface {
	GetName() string
	Broadcast(replicaID int64, name string) error
	Update(data interface{}) error
	SetOnChanged(handler ChangeHandlerInterface)
}
