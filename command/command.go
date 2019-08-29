package command

import "context"

//go:generate mockgen -source command.go -package mock -destination ../mock/command.go

// Command base command interface
type Command interface {
	Command() string
	Reply(ctx context.Context) (string, error)
}
