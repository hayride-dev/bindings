package repository

import (
	"fmt"

	modelrepository "github.com/hayride-dev/bindings/go/internal/gen/hayride/ai/model-repository"
)

// Download downloads a model by its name and returns the path to the model file.
func Download(modelName string) (string, error) {
	modelResult := modelrepository.Download(modelName)
	if modelResult.IsErr() {
		return "", fmt.Errorf("failed to download model %s: %s", modelName, modelResult.Err().Data())
	}

	return *modelResult.OK(), nil
}
