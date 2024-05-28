package schema

import (
	"log"

	"github.com/DimaGitHahahab/ozon-fintech-posts/internal/resolvers"
	"github.com/graphql-go/graphql"
)

func postsField(postType *graphql.Object, resolver *resolvers.Resolver) *graphql.Field {
	return &graphql.Field{
		Type:        graphql.NewList(postType),
		Description: "Get all posts",
		Resolve: func(p graphql.ResolveParams) (any, error) {
			res, err := resolver.GetPosts(p.Context)
			logIfNotNil(err)
			return res, err
		},
	}
}

func postField(postType *graphql.Object, resolver *resolvers.Resolver) *graphql.Field {
	return &graphql.Field{
		Type:        postType,
		Description: "Get post by id",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Resolve: func(p graphql.ResolveParams) (any, error) {
			id, _ := p.Args["id"].(int)
			res, err := resolver.GetPost(p.Context, resolvers.PostArgs{ID: id})
			logIfNotNil(err)
			return res, err
		},
	}
}

func commentsByPostField(commentType *graphql.Object, resolver *resolvers.Resolver) *graphql.Field {
	return &graphql.Field{
		Type:        graphql.NewList(commentType),
		Description: "Get comments by post id",
		Args: graphql.FieldConfigArgument{
			"postId": &graphql.ArgumentConfig{Type: graphql.Int},
			"limit":  &graphql.ArgumentConfig{Type: graphql.Int},
			"offset": &graphql.ArgumentConfig{Type: graphql.Int},
		},
		Resolve: func(p graphql.ResolveParams) (any, error) {
			postId, _ := p.Args["postId"].(int)
			limit, _ := p.Args["limit"].(int)
			offset, _ := p.Args["offset"].(int)
			res, err := resolver.GetCommentsByPost(p.Context, resolvers.GetCommentsArgs{
				PostID: postId,
				Limit:  limit,
				Offset: offset,
			})
			logIfNotNil(err)
			return res, err
		},
	}
}

func commentsByParentField(commentType *graphql.Object, resolver *resolvers.Resolver) *graphql.Field {
	return &graphql.Field{
		Type:        graphql.NewList(commentType),
		Description: "Get comments by parent comment id",
		Args: graphql.FieldConfigArgument{
			"parentId": &graphql.ArgumentConfig{Type: graphql.Int},
			"limit":    &graphql.ArgumentConfig{Type: graphql.Int},
			"offset":   &graphql.ArgumentConfig{Type: graphql.Int},
		},
		Resolve: func(p graphql.ResolveParams) (any, error) {
			parentId, _ := p.Args["parentId"].(int)
			limit, _ := p.Args["limit"].(int)
			offset, _ := p.Args["offset"].(int)
			// making parentId nil if it's 0
			var pointerToParent *int
			if parentId != 0 {
				pointerToParent = &parentId
			}
			res, err := resolver.GetCommentsByParent(p.Context, resolvers.GetCommentsArgs{
				ParentID: pointerToParent,
				Limit:    limit,
				Offset:   offset,
			})
			logIfNotNil(err)
			return res, err
		},
	}
}

func createPostField(postType *graphql.Object, resolver *resolvers.Resolver) *graphql.Field {
	return &graphql.Field{
		Type:        postType,
		Description: "Create new post",
		Args: graphql.FieldConfigArgument{
			"title":    &graphql.ArgumentConfig{Type: graphql.String},
			"content":  &graphql.ArgumentConfig{Type: graphql.String},
			"authorId": &graphql.ArgumentConfig{Type: graphql.Int},
		},
		Resolve: func(p graphql.ResolveParams) (any, error) {
			title, _ := p.Args["title"].(string)
			content, _ := p.Args["content"].(string)
			authorId, _ := p.Args["authorId"].(int)
			res, err := resolver.CreatePost(p.Context, resolvers.CreatePostArgs{
				Title:    title,
				Content:  content,
				AuthorID: authorId,
			})
			logIfNotNil(err)
			return res, err
		},
	}
}

func createCommentField(commentType *graphql.Object, resolver *resolvers.Resolver) *graphql.Field {
	return &graphql.Field{
		Type:        commentType,
		Description: "Create new comment",
		Args: graphql.FieldConfigArgument{
			"postId":   &graphql.ArgumentConfig{Type: graphql.Int},
			"parentId": &graphql.ArgumentConfig{Type: graphql.Int},
			"authorId": &graphql.ArgumentConfig{Type: graphql.Int},
			"content":  &graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: func(p graphql.ResolveParams) (any, error) {
			postId, _ := p.Args["postId"].(int)
			parentId, _ := p.Args["parentId"].(int)
			authorId, _ := p.Args["authorId"].(int)
			content, _ := p.Args["content"].(string)
			// parentId == 0 means that we want to create root comment (null parent id)
			var pointerToParent *int
			if parentId != 0 {
				pointerToParent = &parentId
			}
			res, err := resolver.CreateComment(p.Context, resolvers.CreateCommentArgs{
				PostID:   postId,
				ParentID: pointerToParent,
				AuthorID: authorId,
				Content:  content,
			})
			logIfNotNil(err)
			return res, err
		},
	}
}

func disableCommentsField(resolver *resolvers.Resolver) *graphql.Field {
	return &graphql.Field{
		Type:        graphql.Boolean,
		Description: "Disable new comments for post",
		Args: graphql.FieldConfigArgument{
			"postId":   &graphql.ArgumentConfig{Type: graphql.Int},
			"authorId": &graphql.ArgumentConfig{Type: graphql.Int},
		},
		Resolve: func(p graphql.ResolveParams) (any, error) {
			postId, _ := p.Args["postId"].(int)
			authorId, _ := p.Args["authorId"].(int)
			res, err := resolver.DisableComments(p.Context, resolvers.DisableCommentsArgs{PostID: postId, AuthorId: authorId})
			logIfNotNil(err)
			return res, err
		},
	}
}

func subscribeField(commentType *graphql.Object, resolver *resolvers.Resolver) *graphql.Field {
	return &graphql.Field{
		Type: commentType,
		Resolve: func(p graphql.ResolveParams) (any, error) {
			return p.Source, nil
		},
		Subscribe: func(p graphql.ResolveParams) (any, error) {
			posts, _ := p.Args["posts"].([]int)

			c := make(chan any)

			go func() {
				for {
					select {
					case <-p.Context.Done():
						close(c)
						return
					case c <- resolver.Subscribe(p.Context, c, posts):
					}
				}
			}()

			return c, nil
		},
	}
}

func logIfNotNil(err error) {
	if err != nil {
		log.Println("Error response:", err)
	}
}
