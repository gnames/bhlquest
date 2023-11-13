package storage

type Page struct {
	ID       uint
	ItemID   uint
	FileName string
	Text     string
	Start    uint
	End      uint
}
