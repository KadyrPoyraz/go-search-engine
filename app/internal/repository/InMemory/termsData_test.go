package repository

import (
	"encoding/json"
	"go-search-engine/app/internal/domain"
	"log"
	"os"
	"reflect"
	"testing"
)

const filepath = "index.json"

func setupTest() domain.TermsDataRepository {
	return NewInMemoryTermDataRepository(filepath)
}

func downTest() {
	err := os.RemoveAll(filepath)
	if err != nil {
		log.Fatal(err)
	}
}

func TestInMemoryTermsDataRepository_SaveTermsData(t *testing.T) {
	repo := setupTest()

	termsDataToSave := domain.TermsData{
		FileTermFreq:  domain.FileTermFrequency{
			"Aboba": {
				"Term1": 1,
			},
		},
		FileTermCount: domain.FileTermCount{
			"Aboba": 1,
		},
	}

	err := repo.SaveTermsData(termsDataToSave)
	if err != nil {
		t.Errorf("FAILED: unexpected error: %s", err.Error())
	}

	testFileData, err := os.ReadFile(filepath)
	if err != nil {
		t.Errorf("FAILED: unexpected error: %s", err.Error())
	}

	gottenTermsData := domain.TermsData{}
	 err = json.Unmarshal(testFileData, &gottenTermsData)
	if err != nil {
		t.Errorf("FAILED: unexpected error: %s", err.Error())
	}

	if !reflect.DeepEqual(termsDataToSave, gottenTermsData) {
		t.Errorf("FAILED: expected %v, got %v", termsDataToSave, gottenTermsData)
	}

	downTest()
}

func TestInMemoryTermsDataRepository_UploadTermsData(t *testing.T) {
	repo := setupTest()

	termsDataToSave := domain.TermsData{
		FileTermFreq:  domain.FileTermFrequency{
			"Aboba": {
				"Term1": 1,
			},
		},
		FileTermCount: domain.FileTermCount{
			"Aboba": 1,
		},
	}

	err := repo.SaveTermsData(termsDataToSave)
	if err != nil {
		t.Errorf("FAILED: unexpected error: %s", err.Error())
	}

	gottenTermsData, err := repo.UploadTermsData()

	expectedTermsData := domain.TermsData{
		FileTermFreq:  termsDataToSave.FileTermFreq,
		FileTermCount: termsDataToSave.FileTermCount,
	}

	if !reflect.DeepEqual(expectedTermsData, gottenTermsData) {
		t.Errorf("FAILED: expected %v, got %v", expectedTermsData, gottenTermsData)
	}

	downTest()
}