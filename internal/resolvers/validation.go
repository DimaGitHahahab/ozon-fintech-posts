package resolvers

import "fmt"

const maxLength = 2000

var ErrInvalidComment = fmt.Errorf("comment is too long")

func validateComment(comment string) error {
	if len(comment) <= maxLength {
		return ErrInvalidComment
	}

	return nil
}
