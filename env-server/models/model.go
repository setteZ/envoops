package models

type NodeData struct {
	Id       int64   `db:"ID" json:"id"`
	Time     string  `db:"TIME" json:"time"`
	Quantity string  `db:"QUANTITY" json:"quantity"`
	Value    float64 `db:"VALUE" json:"value"`
}
