package updater

import (
	"context"
	"github.com/jvelo/icecast-monitor/db"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"time"
)

const upsertTrackQuery = `
-- Upsert track
WITH track_cast AS (SELECT id FROM stream WHERE name = $1 LIMIT 1),
     latest_track AS (SELECT * FROM track WHERE cast_id = (SELECT id FROM track_cast) ORDER BY started_at DESC LIMIT 1),
     new_track (title, started_at, ended_at, listeners) AS (SELECT * FROM (VALUES ($2, $3::timestamp, $4::timestamp, $5::int)) v),
     raw (title, started_at, cast_id, ended_at, listeners) AS (
         SELECT title,
                CASE
                    WHEN (SELECT title FROM new_track) = (SELECT title FROM latest_track) THEN
                    	 (SELECT started_at FROM latest_track)
                    ELSE (SELECT started_at FROM new_track)
                    END,
                (SELECT id FROM track_cast),
                ended_at,
                CASE
                    WHEN (SELECT title FROM new_track) = (SELECT title FROM latest_track) THEN
                         GREATEST((SELECT listeners from new_track), COALESCE((SELECT listeners FROM latest_track), -1))
                    ELSE (SELECT listeners FROM new_track)
                    END
         FROM new_track
     )
INSERT
INTO track(title, started_at, cast_id, ended_at, listeners)
SELECT * FROM raw WHERE true
ON CONFLICT (title, started_at, cast_id) DO UPDATE
    SET listeners = (SELECT listeners FROM raw),
        ended_at  = (SELECT ended_at FROM raw)
`

const upsertCastQuery = `
-- Upsert cast
WITH new_cast (name, description, url, updated_at) AS (
    SELECT *
    FROM (VALUES ($1, $2, $3, $4::timestamp)) v
)
INSERT INTO stream (name, description, url, updated_at)
SELECT * FROM new_cast WHERE true
ON CONFLICT (url) DO UPDATE
    SET name        = (SELECT name FROM new_cast),
        description = (SELECT description FROM new_cast),
        updated_at  = (SELECT updated_at FROM new_cast)
`

func upsertTrack(ctx context.Context, record *db.Record) error {
	now := time.Now()
	_, err := boil.GetDB().Exec(upsertTrackQuery, record.Stream.Name, record.Track.Title, now, now, record.Track.Listeners)
	if err != nil {
		return err
	}
	return nil
}

func upsertCast(ctx context.Context, cast *db.Cast) error {
	_, err := boil.GetDB().Exec(upsertCastQuery, cast.Name, cast.Description, cast.URL, time.Now())
	if err != nil {
		return err
	}
	return nil
}
