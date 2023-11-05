package fileredirect

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/mehmetumit/dexus/internal/core/ports"
	"github.com/mehmetumit/dexus/internal/mocks"
)

var (
	// Make paths usable on different oses
	redirectionFilePath                = filepath.ToSlash("./../../../test_data/redirection.yaml")
	faultyFormattedRedirectionFilePath = filepath.ToSlash("./../../../test_data/faulty_formatted_redirection.yaml")
	notExistsRedirectionFilePath       = filepath.ToSlash("./../../../this/path/not/exists")
)

func newTestYamlRedirect(tb testing.TB) (*YamlRedirect, error) {
	tb.Helper()
	l := mocks.NewMockLogger()
	return NewYamlRedirect(l, redirectionFilePath)
}

func TestYamlRedirect_FaultyFormattedFile(t *testing.T) {
	mockLogger := mocks.NewMockLogger()

	_, err := NewYamlRedirect(mockLogger, faultyFormattedRedirectionFilePath)
	if err != ErrUnableToDecodeYamlFile {
		t.Errorf("Expected err %v, got %v", ErrUnableToOpenYamlFile, err)
	}

}
func TestYamlRedirect_NotExistsFile(t *testing.T) {
	mockLogger := mocks.NewMockLogger()
	_, err := NewYamlRedirect(mockLogger, notExistsRedirectionFilePath)
	if err != ErrUnableToOpenYamlFile {
		t.Errorf("Expected err %v, got %v", ErrUnableToOpenYamlFile, err)
	}
}
func TestYamlRedirect_CorrectFile(t *testing.T) {
	mockLogger := mocks.NewMockLogger()
	_, err := NewYamlRedirect(mockLogger, redirectionFilePath)
	if err != nil {
		t.Errorf("Expected err %v, got %v", nil, err)
	}
}
func TestYamlRedirect_NotInitializedStore(t *testing.T) {
	mockLogger := mocks.NewMockLogger()
	hasNotInitializedStore := YamlRedirect{logger: mockLogger}
	ctx := context.Background()
	_, err := hasNotInitializedStore.Get(ctx, "/abc")
	if err != ErrStoreIsNotInitialized {
		t.Errorf("Expected err %v, got %v", ErrStoreIsNotInitialized, err)
	}
}
func TestYamlRedirect_RedirectionFound(t *testing.T) {
	yr, err := newTestYamlRedirect(t)
	ctx := context.Background()
	if err != nil {
		t.Errorf("Expected err %v, got %v", nil, err)
	}
	for k, v := range yr.redirectionStore.RedirectionMap {
		gotRedirect, err := yr.Get(ctx, k)
		if err != nil {
			t.Errorf("Expected err %v, got %v", nil, err)
		}
		if gotRedirect != v {
			t.Errorf("Expected redirect %v, got %v", v, gotRedirect)
		}

	}
}
func TestYamlRedirect_RedirectionNotFound(t *testing.T) {
	yr, err := newTestYamlRedirect(t)
	ctx := context.Background()
	if err != nil {
		t.Errorf("Expected err %v, got %v", nil, err)
	}
	_, err = yr.Get(ctx, "/this/path/not/exists")
	if err != ports.ErrRedirectionNotFound {
		t.Errorf("Expected err %v, got %v", ports.ErrRedirectionNotFound, err)
	}
}
