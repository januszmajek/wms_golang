package stock

type ReportItem struct {
	ProductID   int64  `json:"productId"`
	ArticleCode string `json:"articleCode"`
	Name        string `json:"name"`
	Quantity    int    `json:"quantity"`
}

type InboundRequest struct {
	ProductID int64 `json:"productId"`
	Quantity  int   `json:"quantity"`
}

type InboundResponse struct {
	ProductID     int64 `json:"productId"`
	QuantityAdded int   `json:"quantityAdded"`
}
