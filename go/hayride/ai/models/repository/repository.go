package repository

import (
	"fmt"

	modelrepository "github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/ai/model-repository"
)

var _ ModelRepository = (*modelRepositoryImpl)(nil)

type ModelRepository interface {
	DownloadModel(modelName string) (string, error)
	GetModel(modelName string) (string, error)
	DeleteModel(modelName string) error
	ListModels() ([]string, error)
}

type modelRepositoryImpl struct{}

func New() ModelRepository {
	return &modelRepositoryImpl{}
}

// Download downloads a model by its name and returns the path to the model file.
func (m *modelRepositoryImpl) DownloadModel(modelName string) (string, error) {
	modelResult := modelrepository.DownloadModel(modelName)
	if modelResult.IsErr() {
		return "", fmt.Errorf("failed to download model %s: %s", modelName, modelResult.Err().Data())
	}

	return *modelResult.OK(), nil
}

func (m *modelRepositoryImpl) GetModel(modelName string) (string, error) {
	modelResult := modelrepository.GetModel(modelName)
	if modelResult.IsErr() {
		return "", fmt.Errorf("failed to get model %s: %s", modelName, modelResult.Err().Data())
	}

	return *modelResult.OK(), nil
}

func (m *modelRepositoryImpl) DeleteModel(modelName string) error {
	deleteResult := modelrepository.DeleteModel(modelName)
	if deleteResult.IsErr() {
		return fmt.Errorf("failed to delete model %s: %s", modelName, deleteResult.Err().Data())
	}

	return nil
}

func (m *modelRepositoryImpl) ListModels() ([]string, error) {
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
