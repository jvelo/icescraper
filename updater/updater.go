package updater

import (
	"context"
	"github.com/jvelo/icecast-monitor/model"
	"github.com/jvelo/icecast-monitor/prisma/db"
	log "github.com/sirupsen/logrus"
	"time"
)

func Update(ctx context.Context, stream <-chan *model.Record) error {
	client := db.NewClient()

	if err := client.Prisma.Connect(); err != nil {
		return err
	}

	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case record := <-stream:
			//if err := createOrReplaceCast(ctx, client, record.Cast); err != nil {
			//	log.Errorf("upserting record: %v", err)
			//}
			log.Infof("upserting track %v", record.Track)
			if err := updateTrack(ctx, client, record); err != nil {
				log.Errorf("upserting track: %v", err)
			}
		}
	}

	return nil
}

func updateTrack(ctx context.Context, client *db.PrismaClient, record *model.Record) error {
	now := time.Now().Format(time.RFC3339)
	query := client.Prisma.ExecuteRaw(`
-- Upsert track
WITH track_cast as (SELECT id
           FROM cast
           WHERE name = ?
           LIMIT 1),
     latest AS (SELECT started_at
                FROM track
                WHERE cast_id IN (
                    select id
                    from track_cast
                )
                ORDER BY started_at DESC
                LIMIT 1)
REPLACE
INTO track(title, started_at, cast_id, ended_at, listeners)
VALUES (?,
        COALESCE((SELECT started_at FROM latest), ?),
        (SELECT id FROM track_cast),
        ?,
        ?);
	`, record.Cast.Name, record.Track.Title, now, now, record.Track.Listeners)
	_, err := query.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func createOrReplaceCast(ctx context.Context, client *db.PrismaClient, cast *model.Cast) error {
	query := client.Prisma.ExecuteRaw(`
-- Upsert cast
REPLACE INTO cast(name, description, url, updated_at)
VALUES(?, ?, ?, ?);
	`, cast.Name, cast.Description, cast.URL, time.Now().Format(time.RFC3339))
	_, err := query.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
