package schema

import (
	"github.com/DimaGitHahahab/ozon-fintech-posts/internal/resolvers"
	"github.com/graphql-go/graphql"
)

func PostsField(postType *graphql.Object, resolver *resolvers.Resolver) *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewList(postType),
		Resolve: func(p graphql.ResolveParams) (any, error) {
			return resolver.GetPosts()
		},
	}
}

func PostField(postType *graphql.Object, resolver *resolvers.Resolver) *graphql.Field {
	return &graphql.Field{
		Type: postType,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Resolve: func(p graphql.ResolveParams) (any, error) {
			id, _ := p.Args["id"].(int)
			return resolver.GetPost(p.Context, PostArgs{ID: id})
		},
	}
}

func CommentsByPostField(commentType *graphql.Object, resolver *resolvers.Resolver) *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewList(commentType),
		Args: graphql.FieldConfigArgument{
			"postId": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"limit": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"offset": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Resolve: func(p graphql.ResolveParams) (any, error) {
			postId, _ := p.Args["postId"].(int)
			limit, _ := p.Args["limit"].(int)
			offset, _ := p.Args["offset"].(int)
			return resolver.GetCommentsByPost(p.Context, GetCommentsArgs{
				PostID: postId,
				Limit:  limit,
				Offset: offset,
			})
		},
	}
}

func CommentsByParentField(commentType *graphql.Object, resolver *resolvers.Resolver) *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewList(commentType),
		Args: graphql.FieldConfigArgument{
			"parentId": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"limit": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"offset": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Resolve: func(p graphql.ResolveParams) (any, error) {
			parentId, _ := p.Args["parentId"].(int)
			limit, _ := p.Args["limit"].(int)
			offset, _ := p.Args["offset"].(int)
			return resolver.GetCommentsByParent(p.Context, GetCommentsArgs{
				ParentID: parentId,
				Limit:    limit,
				Offset:   offset,
			})
		},
	}
}

func CreatePostField(postType *graphql.Object, resolver *resolvers.Resolver) *graphql.Field {
	return &graphql.Field{
		Type: postType,
		Args: graphql.FieldConfigArgument{
			"title": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"content": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"authorId": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Resolve: func(p graphql.ResolveParams) (any, error) {
			title, _ := p.Args["title"].(string)
			content, _ := p.Args["content"].(string)
			authorId, _ := p.Args["authorId"].(int)
			return resolver.CreatePost(p.Context, CreatePostArgs{
				Title:    title,
				Content:  content,
				AuthorID: authorId,
			})
		},
	}
}

func CreateCommentField(commentType *graphql.Object, resolver *resolvers.Resolver) *graphql.Field {
	return &graphql.Field{
		Type: commentType,
		Args: graphql.FieldConfigArgument{
			"postId": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"parentId": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"authorId": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"content": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(p graphql.ResolveParams) (any, error) {
			postId, _ := p.Args["postId"].(int)
			parentId, _ := p.Args["parentId"].(int)
			authorId, _ := p.Args["authorId"].(int)
			content, _ := p.Args["content"].(string)
			return resolver.CreateComment(p.Context, CreateCommentArgs{
				PostID:   postId,
				ParentID: parentId,
				AuthorID: authorId,
				Content:  content,
			})
		},
	}
}

func DisableCommentsField(resolver *resolvers.Resolver) *graphql.Field {
	return &graphql.Field{
		Type: graphql.Boolean,
		Args: graphql.FieldConfigArgument{
			"postId": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Resolve: func(p graphql.ResolveParams) (any, error) {
			postId, _ := p.Args["postId"].(int)
			return resolver.DisableComments(p.Context, DisableCommentsArgs{PostID: postId})
		},
	}
}
