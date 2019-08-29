package telegram

// botLogger implements silent logger for bot
type botLogger struct{}

func (l *botLogger) Println(v ...interface{})               {}
func (l *botLogger) Printf(format string, v ...interface{}) {}
