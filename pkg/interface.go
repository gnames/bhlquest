package bhlquest

// BHLQuest provides functionality needed to apply AI to BHL.
type BHLQuest interface {
	// Init bootstraps AI engines providing necessary data and metadata.
	Init() error
}
