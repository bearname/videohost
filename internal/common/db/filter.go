package db

type Page struct {
	Size   int
	Number int
}

type Order string

const (
	Asc  Order = "ASC"
	Desc Order = "DESC"
)

type OrderBy struct {
	FieldName string
	Order     Order
}
