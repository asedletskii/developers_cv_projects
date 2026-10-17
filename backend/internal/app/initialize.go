//go:build !wireinject

package app

func Initialize() (*Application, error) {
	application, _, err := initialize()
	if err != nil {
		return nil, err
	}
	return application, nil
}
