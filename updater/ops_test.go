package updater

import (
	"context"
	"github.com/jvelo/icecast-monitor/model"
	"github.com/jvelo/icecast-monitor/prisma/db"
	"os"
	"testing"
)

var (
	client *db.PrismaClient
	cast   = model.NewCast("TBS", "Test Broadcasting Station", "https://tbs.radio")
)

func setUp() *db.PrismaClient {
	if err := os.Setenv("DATABASE_URL", "file:../prisma/test.db"); err != nil {
		panic(err)
	}

	client = db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		panic(err)
	}

	query1 := client.Prisma.ExecuteRaw("DELETE FROM cast;")
	query2 := client.Prisma.ExecuteRaw("DELETE FROM track;")

	if _, err := query1.Exec(context.Background()); err != nil {
		panic(err)
	}
	if _, err := query2.Exec(context.Background()); err != nil {
		panic(err)
	}

	err := upsertCast(context.Background(), client, cast)
	if err != nil {
		panic(err)
	}

	return client
}

func tearDown() {
	if err := client.Prisma.Disconnect(); err != nil {
		panic(err)
	}
	_ = os.Unsetenv("DATABASE_URL")
}

func TestMain(m *testing.M) {
	setUp()
	result := m.Run()
	tearDown()
	os.Exit(result)
}

func TestUpsertCast(t *testing.T) {
	err := upsertCast(context.Background(), client, cast)
	if err != nil {
		t.Errorf("insert or update cast: %v", err)
	}

	err = upsertCast(context.Background(), client, cast)
	if err != nil {
		t.Errorf("insert or update cast: %v", err)
	}

	casts, err := client.Cast.FindMany().Exec(context.Background())
	if len(casts) != 1 {
		t.Errorf("expected 1 cast, got :%d", len(casts))
	}
}

func TestUpsertTrack(t *testing.T) {
	track := model.NewTrack("Darude – Sandstorm", 42)
	record := &model.Record{
		Cast:  cast,
		Track: track,
	}
	err := upsertTrack(context.Background(), client, record)
	if err != nil {
		t.Errorf("updating track: %v", err)
	}

	err = upsertTrack(context.Background(), client, record)
	if err != nil {
		t.Errorf("updating track: %v", err)
	}
}

func TestUpsertDifferentTracks(t *testing.T) {
	track1 := model.NewTrack("Darude – Sandstorm", 42)
	record1 := &model.Record{
		Cast:  cast,
		Track: track1,
	}
	track2 := model.NewTrack("Dj Fou - Je mets le Waï (Original Version)", 77)
	record2 := &model.Record{
		Cast:  cast,
		Track: track2,
	}
	err := upsertTrack(context.Background(), client, record1)
	if err != nil {
		t.Errorf("upserting track: %v", err)
	}

	err = upsertTrack(context.Background(), client, record2)
	if err != nil {
		t.Errorf("upserting track: %v", err)
	}

	tracksQuery := client.Track.FindMany()
	tracks, err := tracksQuery.Exec(context.Background())

	if err != nil {
		t.Errorf("querying tracks: %v", err)
	}
	if len(tracks) != 2 {
		t.Errorf("expected 2 tracks, got: %v", len(tracks))
	}
}

func TestUpdateTrackListeners(t *testing.T) {
	track := model.NewTrack("Darude – Sandstorm", 42)
	record := &model.Record{
		Cast:  cast,
		Track: track,
	}
	err := upsertTrack(context.Background(), client, record)
	if err != nil {
		t.Errorf("updating track: %v", err)
	}

	track.Listeners = 37
	err = upsertTrack(context.Background(), client, record)
	if err != nil {
		t.Errorf("updating track: %v", err)
	}

	fetched, err := client.Track.FindFirst().Exec(context.Background())
	if err != nil {
		t.Errorf("querying track: %v", err)
	}

	if fetched.Listeners != 42 {
		t.Errorf("expected 42 listeners, got: %v", fetched.Listeners)
	}

	track.Listeners = 51
	err = upsertTrack(context.Background(), client, record)
	if err != nil {
		t.Errorf("updating track: %v", err)
	}

	fetched, err = client.Track.FindFirst().Exec(context.Background())
	if err != nil {
		t.Errorf("querying track: %v", err)
	}

	if fetched.Listeners != 51 {
		t.Errorf("expected 51 listeners, got: %v", fetched.Listeners)
	}
}
