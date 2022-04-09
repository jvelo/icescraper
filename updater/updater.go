package updater

import (
	"context"
	"github.com/jvelo/icecast-monitor/db"
	log "github.com/sirupsen/logrus"
)

func Update(ctx context.Context, stream <-chan *db.Record) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case record := <-stream:
			if err := upsertCast(ctx, record.Cast); err != nil {
				log.Errorf("upserting db: %v", err)
			}
			if err := upsertTrack(ctx, record); err != nil {
				log.Errorf("upserting track: %v", err)
			}
		}
	}
}
