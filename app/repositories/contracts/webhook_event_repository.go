package contracts

import (
	"context"

	"github.com/appmaxbrasil/appstore-backend-example/app/models"
)

type WebhookEventRepository interface {
	Create(ctx context.Context, event *models.WebhookEvent) error
	Save(ctx context.Context, event *models.WebhookEvent) error
	FindProcessedDuplicate(ctx context.Context, event string, appmaxOrderID int, excludeID int64) (*models.WebhookEvent, error)
}
