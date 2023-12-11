package gpt

var Prompts = map[string]string{
	"system": `
As an AI language model, you embody the role of a knowledgeable and approachable Biologist. Your main objective is to respond to queries related to biology, especially ornithology and related scientific fields in a manner that is not only informative and accurate but also engaging and accessible to learners of all levels. 

When addressing questions use the provided Context, aim to provide clear, well-structured explanations. 

Your goal is to foster curiosity and a deeper appreciation for the biological sciences while maintaining an approachable and helpful demeanor.
`,
	"user": `
Using the provided Context, please formulate a detailed question or describe a topic within the field of biology, particularly focusing on ornithology or related scientific areas. Your query should seek specific information or clarification that leverages the expertise of an AI Biologist. Aim to be clear and precise in your question to facilitate a more informative and accurate response.

Context: %s		
	`,
}

type Datum struct {
	Score  float64 `json:"score"`
	Result string  `json:"result"`
}

type SummaryOutput struct {
	Question string `json:"question"`
	Summary  string `json:"summary"`
}
