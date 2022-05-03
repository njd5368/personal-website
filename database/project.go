package database

type Project struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Type        string `json:"type"`
	Image       int64  `json:"image"`

	Languages    []string `json:"languages"`
	Technologies []string `json:"technologies"`

	File []byte `json:"file"`
}
