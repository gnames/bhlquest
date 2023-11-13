package embed

type Chunk struct {
	ID        int
	TitleID   int
	PageID    int
	UUID      string
	Embedding []float32
}
