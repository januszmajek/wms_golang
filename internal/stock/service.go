package stock

import "errors"

var ErrBadQuantity = errors.New("quantity must be bigger than 0")
var ErrProductMissing = errors.New("product not found")

type Service struct{ Repo RepositoryInterface }

func NewService(repo RepositoryInterface) *Service { return &Service{Repo: repo} }

func (s *Service) Receive(req InboundRequest) (InboundResponse, error) {
	if req.Quantity <= 0 {
		return InboundResponse{}, ErrBadQuantity
	}
	exists, err := s.Repo.ProductExists(req.ProductID)
	if err != nil {
		return InboundResponse{}, err
	}
	if !exists {
		return InboundResponse{}, ErrProductMissing
	}
	if err := s.Repo.AddInbound(req.ProductID, req.Quantity); err != nil {
		return InboundResponse{}, err
	}
	return InboundResponse{ProductID: req.ProductID, QuantityAdded: req.Quantity}, nil
}

func (s *Service) Report() ([]ReportItem, error) { return s.Repo.Report() }
