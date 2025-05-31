package object

type Pools struct {
	ID          int    `csv:"id"`
	Name        string `csv:"name"`
	CreatedAt   string `csv:"created_at"`
	UpdatedAt   string `csv:"updated_at"`
	CreatorId   int    `csv:"creator_id"`
	Description string `csv:"description"`
	IsActive    bool   `csv:"is_active"`
	Category    string `csv:"category"`
	PostIds     string `csv:"post_ids"`
}
