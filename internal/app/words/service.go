package words

import (
	"context"
	"strings"
)

type Service interface {
	Create(context context.Context, ownerID string, in CreateWordReq) (string, error)
}

type service struct{ theRepository Repository }

func NewService(repository Repository) Service { return &service{theRepository: repository} }

func (service *service) Create(context context.Context, ownerID string, in CreateWordReq) (string, error) {
	text := strings.TrimSpace(in.Text)

	w := &Word{
		Text:    text,
		OwnerID: ownerID, // always from auth, never from client body
	}
	return service.theRepository.Create(context, w)
}
