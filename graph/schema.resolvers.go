package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"strconv"

	"github.com/jvelo/icecast-monitor/graph/generated"
	"github.com/jvelo/icecast-monitor/graph/model"
	log "github.com/sirupsen/logrus"
)

func (r *queryResolver) Casts(ctx context.Context) ([]*model.Cast, error) {
	models, err := r.Prisma.Cast.FindMany().Exec(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*model.Cast, len(models))
	for i, m := range models {
		description, ok := m.Description()
		if !ok {
			log.Warnf("fetching description for model: %v", m.ID)
		}
		result[i] = &model.Cast{
			ID:          strconv.Itoa(m.ID),
			Name:        m.Name,
			Description: &description,
			URL:         m.URL,
			UpdatedAt:   m.UpdatedAt,
		}
	}
	return result, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
