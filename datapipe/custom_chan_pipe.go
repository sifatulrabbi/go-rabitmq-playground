package datapipe

import (
	"context"
	"encoding/json"
)

type (
	Bytes []byte

	Event struct {
		Name string
		Body Bytes
	}

	Consumer struct {
		Name      string
		EventName string
		Chan      chan Bytes
	}
)

var events = make(chan Event)

func NewCustomDataPipe(ctx context.Context, consumers []*Consumer) {
	for {
		select {
		case e := <-events:
			for i := 0; i < len(consumers); i++ {
				c := consumers[i]

				// if c.EventName != e.Name {
				// 	continue
				// }

				select {
				case c.Chan <- e.Body:
				}

			}
		}
	}
}

func PublishToCustom(name string, body interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	select {
	case events <- Event{Name: name, Body: b}:
	}

	return nil
}
