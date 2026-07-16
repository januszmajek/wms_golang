package stock

import "errors"

var (
	ErrBadQuantity    = errors.New("quantity must be bigger than 0")
	ErrProductMissing = errors.New("product not found")
)

type Service struct{ repo *Repository }

func NewService(repo *Repository) *Service { return &Service{repo: repo} }

func (s *Service) Receive(req InboundRequest) (InboundResponse, error) {
	if req.Quantity <= 0 {
		return InboundResponse{}, ErrBadQuantity
	}
	exists, err := s.repo.ProductExists(req.ProductID)
	if err != nil {
		return InboundResponse{}, err
	}
	if !exists {
		return InboundResponse{}, ErrProductMissing
	}
	if err := s.repo.AddInbound(req.ProductID, req.Quantity); err != nil {
		return InboundResponse{}, err
	}
	return InboundResponse{ProductID: req.ProductID, QuantityAdded: req.Quantity}, nil
}

func (s *Service) Report() ([]ReportItem, error) { return s.repo.Report() }
