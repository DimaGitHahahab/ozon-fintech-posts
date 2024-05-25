package schema

import (
	"github.com/DimaGitHahahab/ozon-fintech-posts/internal/resolvers"
	"github.com/graphql-go/graphql"
)

func NewSchema(resolver *resolvers.Resolver) (graphql.Schema, error) {
	post := PostObject()
	comment := CommentObject()

	rootQuery := getQuery(post, comment, resolver)
	rootMutation := getMutation(post, comment, resolver)

	schemaConfig := graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	}
	return graphql.NewSchema(schemaConfig)
}

func getQuery(postType, commentType *graphql.Object, resolver *resolvers.Resolver) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"posts":            PostsField(postType, resolver),
			"post":             PostField(postType, resolver),
			"commentsByPost":   CommentsByPostField(commentType, resolver),
			"commentsByParent": CommentsByParentField(commentType, resolver),
		},
	})
}

func getMutation(postType, commentType *graphql.Object, resolver *resolvers.Resolver) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"createPost":      CreatePostField(postType, resolver),
			"createComment":   CreateCommentField(commentType, resolver),
			"disableComments": DisableCommentsField(resolver),
		},
	})
}
