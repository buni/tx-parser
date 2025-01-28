package testutils

func ToPtr[T any](v T) *T {
	return &v
}
