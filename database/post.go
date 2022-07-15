package database

type Post struct {
	ID          int64  `json:"id"`
	Type		Type `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Category        string `json:"category"`
	Image       int64  `json:"image"`

	Languages    []string `json:"languages"`
	Technologies []string `json:"technologies"`

	File []byte `json:"file"`
}
