package database

type Project struct {
	ID          int64
	Name        string
	Description string
	Date        string
	Type        string
	Image       []byte

	Languages    []string
	Technologies []string

	File []byte
}
