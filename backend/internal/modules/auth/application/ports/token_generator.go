package ports

type TokenGenerator interface {
	GenerateState() ([]byte, error)
	GenerateBrowserBinding() ([]byte, error)
	GenerateNonce() ([]byte, error)
	GeneratePKCEVerifier() ([]byte, error)
	GenerateSessionToken() ([]byte, error)
	Hash(value []byte) []byte
}
