package bhln

// BHLN takes data from BHLnames.
type BHLN interface {
	ItemIds(offset, limit int, taxa []string) ([]uint, error)
	// Close cleans up database connections
	Close()
}
