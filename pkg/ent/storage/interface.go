package storage

type Storage interface {
	Pages(itemID uint) ([]Page, error)
}
