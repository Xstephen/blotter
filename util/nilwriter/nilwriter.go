package nilwriter

// NilWriter no output
type NilWriter struct{}

func (o NilWriter) Write(p []byte) (int, error) {
	return 0, nil
}
