package data

type Data struct {
	FileTermFreq map[string]map[string]int
	FileTermCount map[string]int
}

func (d *Data) AddFileTermFreqItem(filePath string, term string) {
	if _, ok := d.FileTermFreq[filePath]; !ok {
		d.FileTermFreq[filePath] = make(map[string]int)
	}

	if _, ok := d.FileTermFreq[filePath][term]; ok {
		d.FileTermFreq[filePath][term] += 1
	} else {
		d.FileTermFreq[filePath][term] = 1
	}
}

func (d *Data) AddFileTermCount(filePath string) {
	if _, ok := d.FileTermCount[filePath]; ok {
		d.FileTermCount[filePath] += 1
	} else {
		d.FileTermCount[filePath] = 1
	}
}
