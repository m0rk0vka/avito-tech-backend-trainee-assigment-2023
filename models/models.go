package models

type Users struct {
	ID int64 `json:"id"`
}

type Segments struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Update struct {
	SegmentsToAdd    []string `json:"segments_to_add"`
	SegmentsToDelete []string `json:"segments_to_delete"`
	UserID           int64    `json:"user_id"`
}
