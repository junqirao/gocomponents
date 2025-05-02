package task

import (
	"sync"
)

var (
	notifyChannels = sync.Map{} // id:channel
)

func notify(task *Task) {
	v, loaded := notifyChannels.LoadAndDelete(buildId(task.Type, task.Id))
	if loaded {
		v.(chan *Task) <- task
	}
}

func getNotifyChannel(id string) <-chan *Task {
	ch := make(chan *Task)
	notifyChannels.Store(id, ch)
	return ch
}
