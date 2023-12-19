package bhln

// BHLN takes data from BHLnames.
type BHLN interface {
	// ItemIds retuns back ids of BHL items that need to be embedded.
	ItemIds(offset, limit int, taxa []string) ([]uint, error)

	// Close cleans up database connections
	Close()
}
