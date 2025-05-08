package object

type Post struct {
	ID             int    `csv:"id"`
	UploaderID     int    `csv:"uploader_id"`
	CreatedAt      string `csv:"created_at"`
	MD5            string `csv:"md5"`
	Source         string `csv:"source"`
	Rating         string `csv:"rating"`
	ImageWidth     int    `csv:"image_width"`
	ImageHeight    int    `csv:"image_height"`
	TagString      string `csv:"tag_string"`
	LockedTags     string `csv:"locked_tags"`
	FavCount       int    `csv:"fav_count"`
	FileExt        string `csv:"file_ext"`
	ParentID       int    `csv:"parent_id"`
	ChangeSeq      int    `csv:"change_seq"`
	ApproverID     int    `csv:"approver_id"`
	FileSize       int    `csv:"file_size"`
	CommentCount   int    `csv:"comment_count"`
	Description    string `csv:"description"`
	Duration       int    `csv:"duration"`
	UpdatedAt      string `csv:"updated_at"`
	IsDeleted      bool   `csv:"is_deleted"`
	IsPending      bool   `csv:"is_pending"`
	IsFlagged      bool   `csv:"is_flagged"`
	Score          int    `csv:"score"`
	UpScore        int    `csv:"up_score"`
	DownScore      int    `csv:"down_score"`
	IsRatingLocked bool   `csv:"is_rating_locked"`
	IsStatusLocked bool   `csv:"is_status_locked"`
	IsNoteLocked   bool   `csv:"is_note_locked"`
}
