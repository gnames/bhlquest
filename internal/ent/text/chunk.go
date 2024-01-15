package text

type Chunk struct {
	ID        uint
	ItemID    uint
	Start     uint
	Length    uint
	Text      string
	Embedding []float32
	Score     float64
}
