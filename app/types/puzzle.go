package types

import "encoding/json"

type Puzzle struct {
	ID      int             `json:"id"`
	Title   string          `json:"title"`
	Author  string          `json:"author"`
	RowSize int             `json:"row_size"`
	ColSize int             `json:"col_size"`
	Clues   json.RawMessage `json:"clues"`
}
