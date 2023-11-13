package text

type Text interface {
	TextToChunks(titleID uint) ([]Chunk, error)
}
