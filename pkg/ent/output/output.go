package output

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
	Results []*Result `json:"results"`
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
	ItemID uint `json:"itemId"  example:"226148"`

	// PageID is the ID of the first page of the answer.
	PageID uint `json:"pageId" example:"53469262"`

	// PageIndex is the index of the starting page of the answer.
	PageIndex int `json:"pageIndex" example:"2"`

	// Page is a list of pages in the Item.
	Pages []Page `json:"pages"`

	// Score, generated by AI, indicates the relevance of
	// the result. Higher scores are better.
	Score float64 `json:"score" example:"0.7505834773704542"`

	// CrossScore is generated by Cross-Embeding during comparison
	// of the question with results. It is used for sorting results.
	CrossScore float64 `json:"crossScore,omitempty" example:"0.92353212"`

	// Outlink is the URL pointing to the BHL web page
	// for PageID.
	Outlink string `json:"outlink" example:"https://www.biodiversitylibrary.org/page/53469262"`

	// TextPages is the text from pages of the chunk.
	TextPages []string `json:"text"`

	// TextForSummary is the text of the page, which is used for summary.
	TextForSummary string `json:"-"`

	// Reference is the string representation of the BHL reference.
	Reference string `json:"reference,omitempty"`

	// Language is the main language of the item's title.
	Language string `json:"language,omitempty"`

	// OutlinkTitleDOI is the DOI of the item's title.
	OutlinkTitleDOI string `json:"outlinkTitleDOI,omitempty"`
}

type Page struct {
	// ID is the ID of the page.
	ID uint `json:"id" example:"53469262"`

	// PageNum is the page number of the page in the item.
	// If it is not given, no page number is available.
	PageNum string `json:"pageNum" example:"2"`
}
