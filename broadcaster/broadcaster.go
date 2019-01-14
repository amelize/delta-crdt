package broadcaster

import "time"

type BroadcastHandler interface {
	OnBroadcast(id string, name string, data []byte) error
}

func New(id string) *Broadcaster {
	return &Broadcaster{
		id: id,
	}
}

type Broadcaster struct {
	id string
}

func (b *Broadcaster) broadcast() {

}

func (b *Broadcaster) loop() {
	// TODO: stop
	ticker := time.NewTicker(time.Millisecond * 500)

	for {
		select {
		case <-ticker.C:
			b.broadcast()

		}
	}
}
