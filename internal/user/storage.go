package user

import (
	"context"
)

type Repository interface {
	JoinGroup(ctx context.Context) error
	LeaveGroup(ctx context.Context) error
}
