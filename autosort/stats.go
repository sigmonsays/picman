package autosort

type Stats struct {
	Processed      int `json:"processed"`
	Copied         int `json:"copied"`
	FilesProcessed int `json:"files_processed"`
	DirsProcessed  int `json:"dirs_processed"`
}
