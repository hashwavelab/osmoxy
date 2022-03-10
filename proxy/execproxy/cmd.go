package execproxy

type Command interface {
	Execute() ([]byte, error)
}
