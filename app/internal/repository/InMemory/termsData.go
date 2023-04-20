package repository

import (
	"encoding/json"
	"go-search-engine/app/internal/domain"
	"os"
)

type inMemoryTermsDataRepository struct {
	filepath string
}

func NewInMemoryTermDataRepository(filepath string) domain.TermsDataRepository {
	return &inMemoryTermsDataRepository{
		filepath: filepath,
	}
}

func (r *inMemoryTermsDataRepository) SaveTermsData(termsData domain.TermsData) error {
	data, err := json.Marshal(termsData)
	if err != nil {
		return nil
	}

	err = os.WriteFile(r.filepath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (r *inMemoryTermsDataRepository) UploadTermsData() (domain.TermsData, error) {
	fileData, err := os.ReadFile(r.filepath)
	if err != nil {
		return domain.TermsData{}, err
	}

	uploadedTermsData := domain.TermsData{}
	err = json.Unmarshal(fileData, &uploadedTermsData)
	if err != nil {
		return domain.TermsData{}, err
	}

	return uploadedTermsData, nil
}


