package datamanager

import (
	"encoding/json"
	"io"
	"os"

	"github.com/terrachest/chest/internal/modules"
)

const (
	tmpFileName  = "tmp_registry_data"
	dataFileName = "/modules.json"
)

type DataManager interface {
	Load(ms *modules.Modules) error
	Persist(ms modules.Modules) error
}

type Config struct {
	DataDir string
}

type dataManager struct {
	config *Config
}

func New(c *Config) DataManager {
	return &dataManager{
		config: c,
	}
}

func (manager *dataManager) Load(ms *modules.Modules) error {
	data, err := manager.readFile(ms)

	if err != nil {
		return err
	}

	return json.Unmarshal(data, ms)
}

func (manager *dataManager) Persist(ms modules.Modules) error {
	data, err := json.Marshal(ms)

	if err != nil {
		return err
	}

	return manager.writeFile(data)
}

func (manager *dataManager) readFile(ms *modules.Modules) ([]byte, error) {
	file, err := os.Open(manager.config.DataDir + dataFileName)
	if err != nil {
		if os.IsNotExist(err) {
			data, err := json.Marshal(ms)

			if err != nil {
				return nil, err
			}

			err = manager.writeFile(data)

			if err != nil {
				return nil, err
			}
			return nil, nil
		}
		return nil, err
	}
	return io.ReadAll(file)
}

func (manager *dataManager) writeFile(data []byte) error {
	tmpFile, err := os.CreateTemp(manager.config.DataDir, tmpFileName)

	if err != nil {
		return err
	}

	_, err = tmpFile.Write(data)

	if err != nil {
		return err
	}

	err = tmpFile.Sync()

	if err != nil {
		return err
	}

	return os.Rename(tmpFile.Name(), manager.config.DataDir+dataFileName)
}
