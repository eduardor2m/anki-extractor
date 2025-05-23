package models

type Note struct {
	ID     int64    `json:"id"`
	GUID   string   `json:"guid"`
	Model  int64    `json:"model"`
	Fields []string `json:"fields"`
	Tags   []string `json:"tags"`
}

type Card struct {
	ID          int64 `json:"id"`
	NoteID      int64 `json:"note_id"`
	DeckID      int64 `json:"deck_id"`
	Interval    int   `json:"interval"`
	Due         int64 `json:"due"`
	TemplateOrd int   `json:"template_order"`
}

type Deck struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Flashcard struct {
	Note     Note   `json:"note"`
	Card     Card   `json:"card"`
	DeckName string `json:"deck_name"`
}
