package tabulator


type Tabulator interface {
	Tabulate([]interface{}) error
}
