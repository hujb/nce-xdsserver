package adapter

import "github.com/nce/nce-xdsserver/common/event"

type ResourceWatcherAdapter interface {
	Notify(event *event.Event)
}
