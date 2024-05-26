package resolvers

type PostArgs struct {
	ID int `json:"id"`
}

type CreatePostArgs struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	AuthorID int    `json:"authorId"`
}

type CreateCommentArgs struct {
	PostID   int    `json:"postId"`
	ParentID int    `json:"parentId"`
	AuthorID int    `json:"authorId"`
	Content  string `json:"content"`
}

type GetCommentsArgs struct {
	PostID   int `json:"postId"`
	ParentID int `json:"parentId"`
	Limit    int `json:"limit"`
	Offset   int `json:"offset"`
}

type DisableCommentsArgs struct {
	PostID   int `json:"postId"`
	AuthorId int `json:"authorId"`
}
