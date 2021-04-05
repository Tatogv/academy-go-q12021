package entities

type Response struct {
	count    int
	next     string
	previous string
	Results  []Berry `json:"results,omitempty"`
}

type Berry struct {
	Name string `json:"name,omitempty"`
	Url  string `json:"url,omitempty"`
}
