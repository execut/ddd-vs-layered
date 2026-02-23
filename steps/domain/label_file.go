package domain

type LabelFile struct {
	Path string
}

func NewLabelFile(path string) (LabelFile, error) {
	return LabelFile{Path: path}, nil
}
