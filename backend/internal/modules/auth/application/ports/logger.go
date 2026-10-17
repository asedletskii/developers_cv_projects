package ports

// WarningLogger is the application-facing logging capability needed for
// best-effort operations. Infrastructure-specific logging fields and levels
// remain outside the application layer.
type WarningLogger interface {
	Warn(message string, err error)
}
