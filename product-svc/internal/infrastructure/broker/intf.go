package broker

// EventRouter interface
type EventRouter interface {
	RegisterHandlers()
	Run() error
	GracefulShutdown() error
}
