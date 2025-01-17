package Repository

type TransactionEnder interface {
	End(err error) error
}

type TransactionUnit[T any] interface {
	Begin() (T, error)
	TransactionEnder
}
