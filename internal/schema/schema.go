package schema

import (
	"github.com/DimaGitHahahab/ozon-fintech-posts/internal/resolvers"
	"github.com/graphql-go/graphql"
)

func NewSchema(resolver *resolvers.Resolver) (graphql.Schema, error) {
	post := postObject()
	comment := commentObject()

	rootQuery := query(post, comment, resolver)
	rootMutation := mutation(post, comment, resolver)

	schemaConfig := graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	}
	return graphql.NewSchema(schemaConfig)
}

func query(postType, commentType *graphql.Object, resolver *resolvers.Resolver) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"posts":            postsField(postType, resolver),
			"post":             postField(postType, resolver),
			"commentsByPost":   commentsByPostField(commentType, resolver),
			"commentsByParent": commentsByParentField(commentType, resolver),
		},
	})
}

func mutation(postType, commentType *graphql.Object, resolver *resolvers.Resolver) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"createPost":      createPostField(postType, resolver),
			"createComment":   createCommentField(commentType, resolver),
			"disableComments": disableCommentsField(resolver),
		},
	})
}
