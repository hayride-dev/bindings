package repository

import (
	"fmt"

	modelrepository "github.com/hayride-dev/bindings/go/internal/gen/hayride/ai/model-repository"
)

// Download downloads a model by its name and returns the path to the model file.
func DownloadModel(modelName string) (string, error) {
	modelResult := modelrepository.DownloadModel(modelName)
	if modelResult.IsErr() {
		return "", fmt.Errorf("failed to download model %s: %s", modelName, modelResult.Err().Data())
	}

	return *modelResult.OK(), nil
}

func GetModel(modelName string) (string, error) {
	modelResult := modelrepository.GetModel(modelName)
	if modelResult.IsErr() {
		return "", fmt.Errorf("failed to get model %s: %s", modelName, modelResult.Err().Data())
	}

	return *modelResult.OK(), nil
}

func DeleteModel(modelName string) error {
	deleteResult := modelrepository.DeleteModel(modelName)
	if deleteResult.IsErr() {
		return fmt.Errorf("failed to delete model %s: %s", modelName, deleteResult.Err().Data())
	}

	return nil
}

func ListModels() ([]string, error) {
	modelsResult := modelrepository.ListModels()
	if modelsResult.IsErr() {
		return nil, fmt.Errorf("failed to list models: %s", modelsResult.Err().Data())
	}

	modelsList := modelsResult.OK()
	if modelsList == nil {
		return nil, fmt.Errorf("no models found")
	}

	return modelsList.Slice(), nil
}
