package logtracer

type writer struct {
	ch chan byte
}

func newWriter() *writer {
	return &writer{make(chan byte)}
}

func (w *writer) Chan() <-chan byte {
	return w.ch
}

func (w *writer) Write(p []byte) (int, error) {
	n := 0
	for _, b := range p {
		w.ch <- b
		n++
	}
	return n, nil
}

func (w *writer) Close() error {
	close(w.ch)
	return nil
}
