package event

type Event struct {
	eType EventType
}

var SERVICE_UPDATE_EVENT = &Event{eType: SERVICE}

func (e *Event) GetType() EventType {
	return e.eType
}
