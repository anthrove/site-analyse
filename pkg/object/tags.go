package object

type Tag struct {
	ID        int    `csv:"id"`
	Name      string `csv:"name"`
	Category  int    `csv:"category"`
	PostCount int    `csv:"post_count"`
}
