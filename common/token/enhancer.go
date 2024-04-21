package token

type Enhancer interface {
	Sign(claims *Claims) (string, error)
	Verify(value string) (*Claims, error)
}
