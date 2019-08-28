package handler

import "context"

//go:generate mockgen -source sender.go -package mock -destination ../mock/sender.go

// Sender base sender interface
type Sender interface {
	Send(ctx context.Context, target string, text string) error
}
