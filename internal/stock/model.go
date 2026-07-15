package stock

type ReportItem struct {
	ProductID   int64  `json:"product_id"`
	ArticleCode string `json:"article_code"`
	Name        string `json:"name"`
	Quantity    int    `json:"quantity"`
}

type InboundRequest struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

type InboundResponse struct {
	ProductID     int64 `json:"product_id"`
	QuantityAdded int   `json:"quantity_added"`
}
