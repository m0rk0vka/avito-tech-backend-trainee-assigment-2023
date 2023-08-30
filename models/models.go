package models

type Users struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type Segments struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Relations struct {
	ID        int64 `json:"id"`
	UserID    int64 `json:"user_id"`
	SegmentID int64 `json:"segment_id"`
}
