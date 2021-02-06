package updater

import (
	"context"
	"github.com/jvelo/icecast-monitor/model"
	"github.com/jvelo/icecast-monitor/prisma/db"
	"time"
)

const upsertTrackQuery = `
-- Upsert track
WITH fresh_track (title, started_at, ended_at, listeners) AS (
    SELECT *
    FROM (VALUES (?, ?, ?, ?))
),
     track_cast as (SELECT id FROM cast WHERE name = ? LIMIT 1),
     latest_track as (select * from track ORDER by started_at desc limit 1)
INSERT
INTO track(title, started_at, cast_id, ended_at, listeners)
VALUES ((SELECT title from fresh_track),
        (
            CASE WHEN (select title from fresh_track) = (select title from latest_track) THEN 
                (select started_at from latest_track)
                ELSE (select started_at from fresh_track) END
        ),
        (SELECT id from track_cast),
        (SELECT ended_at from fresh_track),
        (SELECT listeners from fresh_track))
ON CONFLICT (title, started_at, cast_id) DO UPDATE
    SET listeners = CASE WHEN (select title from fresh_track) = (select title from latest_track)  THEN
        	MAX((SELECT listeners from fresh_track),
                        COALESCE((SELECT listeners FROM latest_track), -1))
        ELSE (SELECT listeners FROM fresh_track) END,
        ended_at  = (SELECT ended_at from fresh_track)
`

const upsertCastQuery = `
-- Upsert cast
WITH fresh_cast (name, description, url, updated_at) AS (
    SELECT * FROM (VALUES (?, ?, ?, ?))
)
INSERT
INTO cast(name, description, url, updated_at)
VALUES ((SELECT name from fresh_cast),
        (SELECT description from fresh_cast),
        (SELECT url from fresh_cast),
        (SELECT updated_at from fresh_cast))
ON CONFLICT (url) DO UPDATE
    SET name = (SELECT name from fresh_cast),
        description = (SELECT description from fresh_cast),
		updated_at = (SELECT updated_at from fresh_cast)
`

func updateTrack(ctx context.Context, client *db.PrismaClient, record *model.Record) error {
	now := time.Now().Format(time.RFC3339)
	query := client.Prisma.ExecuteRaw(upsertTrackQuery, record.Track.Title, now, now, record.Track.Listeners, record.Cast.Name)
	_, err := query.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func createOrReplaceCast(ctx context.Context, client *db.PrismaClient, cast *model.Cast) error {
	query := client.Prisma.ExecuteRaw(upsertCastQuery, cast.Name, cast.Description, cast.URL, time.Now().Format(time.RFC3339))
	_, err := query.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
