package cased

import "time"

type MockPublisher struct {
	Events  []AuditEvent
	silence bool
}

func NewMockPublisher() (*MockPublisher, func()) {
	mp := &MockPublisher{silence: false}
	cp := CurrentPublisher()
	SetPublisher(mp)
	closeFunc := func() {
		SetPublisher(cp)
	}

	return mp, closeFunc
}

func NewSilencedMockPublisher() (*MockPublisher, func()) {
	mp := &MockPublisher{silence: true}
	cp := CurrentPublisher()
	SetPublisher(mp)
	closeFunc := func() {
		SetPublisher(cp)
	}

	return mp, closeFunc
}

func (mp MockPublisher) Options() PublisherOptions {
	return PublisherOptions{
		Silence: mp.silence,
	}
}

func (mp MockPublisher) Flush(_ time.Duration) bool {
	return true
}

func (mp *MockPublisher) Publish(event AuditEvent) error {
	mp.Events = append(mp.Events, event)

	return nil
}
