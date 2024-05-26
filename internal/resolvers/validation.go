package resolvers

const maxLength = 2000

func validateComment(comment string) bool {
	return len(comment) <= maxLength
}
