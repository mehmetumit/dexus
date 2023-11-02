package mocks

import (
	"context"

	"github.com/mehmetumit/dexus/internal/core/ports"
)

type MockRedirectionRepo struct {
	redirectionMap MockRedirectionMap
}

type MockRedirectionMap map[string]string

func NewMockRedirectionRepo(p ...MockRedirectionMap) *MockRedirectionRepo {
	return &MockRedirectionRepo{}

}
func (mr *MockRedirectionRepo) Get(ctx context.Context, from string) (string, error) {
	to, ok := mr.redirectionMap[from]
	if !ok {
		return "", ports.ErrRedirectionNotFound
	}
	return to, nil
}
func (mr *MockRedirectionRepo) SetMockRedirectionMap(m MockRedirectionMap) {
	mr.redirectionMap = m
}
