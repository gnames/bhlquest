package answer

// Answer is a struct that represents the response to a query.
// @Description Answer holds the metadata and results for a
// @Description query response, containing the pages from BHL
// @Description that answer a given question.
type Answer struct {
	// Meta contains metadata about the request.
	Meta `json:"metadata"`

	// Summary contains an answer from LLM created from the
	// content of the results.
	Summary string `json:"summary"`

	// Results is a list of pages containing the answers.
	Results []Result `json:"results"`
}

// Meta contains metadata about the request.
// @Description Meta includes information such as the
// @Description time taken to process the query.
type Meta struct {
	// Question asked by user.
	Question string `json:"question" example:"What are ecological niches for Indigo Bunting?"`

	// MaxResultsNum is the maximum number of returned results.
	MaxResultsNum int `json:"maxResultsNum" example:"10"`

	// ScoreThreshold determines the smallest score which is
	// still considered for results.
	ScoreThreshold float64 `json:"scoreThreshold" example:"0.4"`

	// QueryTime is the duration taken to process the query.
	QueryTime float64 `json:"queryTime" example:"0.911422974"`

	// Version of BHLQuest
	Version string `json:"version" example:"v0.0.3"`
}

// Result represents a specific answer found in BHL.
// @Description Result holds information about a BHL page or
// @Description range of pages that contain answers to
// @Description a given question.
type Result struct {
	ChunkID uint `json:"chunkId" example:"2980234"`
	// ItemID is the ID of a BHL Item, such as a book
	// or journal volume.
	ItemID uint `json:"itemId" example:"226148"`

	// PageIDStart is the ID of the starting page of the answer.
	PageIDStart uint `json:"pageStart" example:"53469262"`

	// PageIDEnd is the ID of the ending page of the answer.
	// It's the same as PageIDStart if the answer is on one page.
	PageIDEnd uint `json:"pageEnd" example:"53469262"`

	// Score, generated by AI, indicates the relevance of
	// the result. Higher scores are better.
	Score float64 `json:"score" example:"0.7505834773704542"`

	// CrossScore is generated by Cross-Embeding during comparison
	// of the question with results. It is used for sorting results.
	CrossScore float64 `json:"crossScore" example:"0.02353212"`

	// Outlink is the URL pointing to the BHL website
	// for PageIDStart.
	Outlink string `json:"outlink" example:"https://www.biodiversitylibrary.org/page/53469262"`

	// Text respresents the actual string that was used for matching by AI.
	Text string `json:"text,omitempty"`

	// TextExt contains more text that Text, allowing to create
	// a better summary.
	TextExt string `json:"-"`
}
