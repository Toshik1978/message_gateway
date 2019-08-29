package handler

import "context"

//go:generate mockgen -source receiver.go -package mock -destination ../mock/receiver.go

// Receiver base receiver interface
type Receiver interface {
	Receive(ctx context.Context)
	Stop()
}
