package token

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

const randomTokenSize = 32

type Generator struct{}

func NewGenerator() *Generator { return &Generator{} }

func (g *Generator) GenerateState() ([]byte, error)        { return g.generate() }
func (g *Generator) GenerateNonce() ([]byte, error)        { return g.generate() }
func (g *Generator) GeneratePKCEVerifier() ([]byte, error) { return g.generate() }
func (g *Generator) GenerateSessionToken() ([]byte, error) { return g.generate() }

func (g *Generator) Hash(value []byte) []byte {
	hash := sha256.Sum256(value)
	result := make([]byte, len(hash))
	copy(result, hash[:])
	return result
}

func (g *Generator) generate() ([]byte, error) {
	value := make([]byte, randomTokenSize)
	if _, err := rand.Read(value); err != nil {
		return nil, fmt.Errorf("generate random auth value: %w", err)
	}
	return value, nil
}
