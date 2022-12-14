package datamanager

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"privateterraformregistry/internal/modules"
)

type DataManager interface {
	Load() error
	Save() error
}

type Config struct {
	DataDir string
}

type dataManager struct {
	modules *modules.Modules
	config  *Config
}

func New(c *Config, ms *modules.Modules) DataManager {
	return &dataManager{
		config:  c,
		modules: ms,
	}
}

func (manager *dataManager) Load() error {
	data, err := manager.readFile()

	if err != nil {
		return err
	}

	if data == nil {
		return nil
	}

	if data != nil {
		err = json.Unmarshal(data, &manager.modules)
	}

	if err != nil {
		return err
	}

	return nil
}

func (manager *dataManager) Save() error {
	data, err := json.Marshal(manager.modules)

	if err != nil {
		return err
	}

	return manager.writeFile(data)
}

func (manager *dataManager) readFile() ([]byte, error) {
	file, err := os.Open(manager.config.DataDir + "/modules.json")
	if err != nil {
		if os.IsNotExist(err) {
			data, err := json.Marshal(manager.modules)

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
	return ioutil.ReadAll(file)
}

func (manager *dataManager) writeFile(data []byte) error {
	tmpFile, err := ioutil.TempFile(manager.config.DataDir, "tmp_registry_data")

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

	return os.Rename(tmpFile.Name(), manager.config.DataDir+"/modules.json")
}
