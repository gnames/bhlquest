package answer

type Answer struct {
	Meta
	Results []Result
}

type Meta struct {
	QueryTime float64 `json:"queryTime"`
}

type Result struct {
	ItemID      uint    `json:"itemId"`
	PageIDStart uint    `json:"pageStart"`
	PageIDEnd   uint    `json:"pageEnd"`
	Score       float64 `json:"score"`
	Outlink     string  `json:"outlink"`
}
