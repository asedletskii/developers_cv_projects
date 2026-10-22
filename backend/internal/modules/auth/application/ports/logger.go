package ports

type WarningLogger interface {
	Warn(message string, err error)
}
