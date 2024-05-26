package schema

import "github.com/graphql-go/graphql"

// postObject is a GraphQL object for domain.Post
func postObject() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Post",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"content": &graphql.Field{
				Type: graphql.String,
			},
			"authorId": &graphql.Field{
				Type: graphql.Int,
			},
			"createdAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"commentsDisabled": &graphql.Field{
				Type: graphql.Boolean,
			},
		},
	})
}

// commentObject is a GraphQL object for domain.Comment
func commentObject() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Comment",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"postId": &graphql.Field{
				Type: graphql.Int,
			},
			"parentId": &graphql.Field{
				Type: graphql.Int,
			},
			"authorId": &graphql.Field{
				Type: graphql.Int,
			},
			"content": &graphql.Field{
				Type: graphql.String,
			},
			"createdAt": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	})
}
