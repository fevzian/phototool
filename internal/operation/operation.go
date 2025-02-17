package operation

type Operation interface {
	Execute() error
	GetName() string
}
