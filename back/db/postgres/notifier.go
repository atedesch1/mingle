package db

import (
	"log"
	"time"

	"github.com/lib/pq"
)

const (
	minReconnectInterval = 10 * time.Second
	maxReconnectInterval = time.Minute
)

type notifier struct {
	listener *pq.Listener
	failed   chan error
}

func newNotifier(dsn, channelName string) (*notifier, error) {
	n := &notifier{failed: make(chan error, 2)}

	listener := pq.NewListener(
		dsn,
		minReconnectInterval,
		maxReconnectInterval,
		n.logListener)

	if err := listener.Listen(channelName); err != nil {
		listener.Close()
		return nil, err
	}

	n.listener = listener
	return n, nil
}

func (n *notifier) logListener(event pq.ListenerEventType, err error) {
	if err != nil {
		log.Printf("listener error: %s\n", err)
	}
	if event == pq.ListenerEventConnectionAttemptFailed {
		n.failed <- err
	}
}

func (n *notifier) fetch(data chan<- []byte) error {
	var fetchCounter uint64
	for {
		select {
		case e := <-n.listener.Notify:
			if e == nil {
				continue
			}
			fetchCounter++
			data <- []byte(e.Extra)
		case err := <-n.failed:
			return err
		case <-time.After(time.Minute):
			go n.listener.Ping()
		}
	}
}

func (n *notifier) close() error {
	close(n.failed)
	return n.listener.Close()
}
