package domain

type TermsDataRepository interface {
	SaveTermsData(TermsData) error
	UploadTermsData() (TermsData, error)
}

type FileTermFrequency map[string]map[Term]int
type FileTermCount map[string]int

type TermsData struct {
	FileTermFreq  FileTermFrequency
	FileTermCount FileTermCount
}
