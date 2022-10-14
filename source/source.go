package source

type Source interface {
	Read() (map[string]interface{}, error)
}
