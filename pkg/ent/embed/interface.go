package embed

type Embed interface {
	Init() error
	Populate(itemIDs []uint) error
}
