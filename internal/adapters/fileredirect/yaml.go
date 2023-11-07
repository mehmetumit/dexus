package fileredirect

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/mehmetumit/dexus/internal/core/ports"
	"gopkg.in/yaml.v3"
)

var (
	ErrStoreIsNotInitialized  = errors.New("store is not initialized")
	ErrUnableToOpenYamlFile   = errors.New("unable to open the file")
	ErrUnableToDecodeYamlFile = errors.New("unable to decode yaml file")
)

type RedirectionStore struct {
	RedirectionMap map[string]string `yaml:"redirections"`
}
type YamlRedirect struct {
	sync.RWMutex
	redirectionStore *RedirectionStore
	logger           ports.Logger
}

func initStore(logger ports.Logger, filePath string) (*RedirectionStore, error) {
	f, err := os.Open(filepath.FromSlash(filePath))
	defer func() {
		err := f.Close()
		if err != nil {
			logger.Error("unable to close file:", err)
		}
	}()
	if err != nil {
		logger.Error(ErrUnableToOpenYamlFile, err)
		return nil, ErrUnableToOpenYamlFile
	}
	store := RedirectionStore{
		RedirectionMap: make(map[string]string),
	}
	if err := yaml.NewDecoder(f).Decode(&store); err != nil {
		logger.Error(ErrUnableToDecodeYamlFile, err)
		return nil, ErrUnableToDecodeYamlFile
	}
	return &store, nil
}
func NewYamlRedirect(logger ports.Logger, filePath string) (*YamlRedirect, error) {
	store, err := initStore(logger, filePath)
	if err != nil {
		logger.Error("failed to create yaml redirect")
	}
	return &YamlRedirect{
		redirectionStore: store,
	}, err
}
func (yr *YamlRedirect) Get(ctx context.Context, from string) (string, error) {
	yr.RLock()
	defer yr.RUnlock()
	if yr.redirectionStore == nil {
		return "", ErrStoreIsNotInitialized
	}
	to, ok := yr.redirectionStore.RedirectionMap[from]
	if !ok {
		return "", ports.ErrRedirectionNotFound
	}
	return to, nil
}
