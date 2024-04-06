package personal_cards

import (
	"context"
)

type Repository interface {
	ShowAllPersonalCards(ctx context.Context) (pc []PersonalCard, err error)
}
