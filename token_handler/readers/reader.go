package readers

type Reader interface {
	Read(path string) (map[string]int, error)
}
