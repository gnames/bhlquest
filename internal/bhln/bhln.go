package bhln

type Item struct {
	ID    int
	Title string
}

type Page struct {
	ID     int
	Number int
	ItemID int
}

type Chunk struct {
	ID     int
	UUID   string
	Text   string
	Vector []float32
}
