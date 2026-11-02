package query

const CouldNotFindAnswer = "Could not find the answer to the question."

type AskRequest struct {
	Question string `json:"question"`
}

type AskResponse struct {
	Answer string `json:"answer"`
}
