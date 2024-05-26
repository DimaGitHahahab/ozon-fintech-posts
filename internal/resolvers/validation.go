package resolvers

import "fmt"

const maxLength = 2000

var (
	ErrInvalidComment        = fmt.Errorf("comment is too long")
	ErrNotPositiveID         = fmt.Errorf("ID must be positive")
	ErrInvalidPaginationArgs = fmt.Errorf("invalid pagination args")
)

func validateComment(comment string) error {
	if len(comment) >= maxLength {
		return ErrInvalidComment
	}

	return nil
}

func validateID(IDs ...int) error {
	for _, id := range IDs {
		if id < 1 {
			return ErrNotPositiveID
		}
	}
	return nil
}

func validatePaginationArgs(limit, offset int) error {
	if limit < 1 || offset < 0 {
		return ErrInvalidPaginationArgs
	}
	return nil
}
