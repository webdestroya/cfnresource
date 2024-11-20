package encoding_test

type StringLike string

const StringLikeValue StringLike = `StrLike`

func Ptrize[T any](v T) *T {
	return &v
}
