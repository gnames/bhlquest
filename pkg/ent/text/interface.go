package text

type Text interface {
	TextToChunks(titleID uint) ([]Chunk, error)
	ChunkText(Chunk) string
}
