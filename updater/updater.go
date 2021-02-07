package updater

import (
	"context"
	"github.com/jvelo/icecast-monitor/model"
	"github.com/jvelo/icecast-monitor/prisma/db"
	log "github.com/sirupsen/logrus"
)

func Update(ctx context.Context, stream <-chan *model.Record) error {
	client := db.NewClient()

	if err := client.Prisma.Connect(); err != nil {
		return err
	}

	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			log.Fatalf("disconnecting: %v", err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case record := <-stream:
			if err := upsertCast(ctx, client, record.Cast); err != nil {
				log.Errorf("upserting record: %v", err)
			}
			if err := upsertTrack(ctx, client, record); err != nil {
				log.Errorf("upserting track: %v", err)
			}
		}
	}
}
