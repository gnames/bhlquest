package embed

type Embed interface {
	Init() error
	Populate(itemIDs []int) error
}
