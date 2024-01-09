package challenge

// Challenge interface to use.
type Challenge interface {
	// Request must always be a base64 string to let pass via different protocols.
	Request() ([]byte, error)
	Fulfil(request []byte) ([]byte, error)
	Check(result []byte) (bool, error)
}
