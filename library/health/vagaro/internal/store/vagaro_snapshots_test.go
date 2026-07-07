// Copyright 2026 Trevin Chow and contributors. Licensed under Apache-2.0. See LICENSE.

package store

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceSnapshots(t *testing.T) {
	ctx := context.Background()
	s := newVagaroTestStore(t)

	require.NoError(t, s.InsertServiceSnapshot(ctx, "93458", "2026-07-01T00:00:00Z", []ServiceRecord{
		{ServiceID: "1", Title: "Cut", PriceCents: 5200},
	}))
	require.NoError(t, s.InsertServiceSnapshot(ctx, "93458", "2026-07-02T00:00:00Z", []ServiceRecord{
		{ServiceID: "1", Title: "Cut", PriceCents: 5500},
		{ServiceID: "2", Title: "Shave", PriceCents: 3000},
	}))

	times, err := s.RecentSnapshotTimes(ctx, "93458", 2)
	require.NoError(t, err)
	require.Len(t, times, 2)
	// Newest first.
	assert.Equal(t, "2026-07-02T00:00:00Z", times[0])
	assert.Equal(t, "2026-07-01T00:00:00Z", times[1])

	newer, err := s.SnapshotServices(ctx, "93458", times[0])
	require.NoError(t, err)
	require.Len(t, newer, 2)
	assert.Equal(t, "1", newer[0].ServiceID)
	assert.Equal(t, 5500, newer[0].PriceCents)

	older, err := s.SnapshotServices(ctx, "93458", times[1])
	require.NoError(t, err)
	require.Len(t, older, 1)
}

func TestWatchBaseline(t *testing.T) {
	ctx := context.Background()
	s := newVagaroTestStore(t)

	_, found, err := s.GetWatchBaseline(ctx, "centralbarber", "34098477")
	require.NoError(t, err)
	assert.False(t, found)

	require.NoError(t, s.UpsertWatchBaseline(ctx, "centralbarber", "34098477", "2026-07-24T10:00:00Z", "2026-07-30"))
	b, found, err := s.GetWatchBaseline(ctx, "centralbarber", "34098477")
	require.NoError(t, err)
	require.True(t, found)
	assert.Equal(t, "2026-07-24T10:00:00Z", b.NextAvailable)
	assert.Equal(t, "2026-07-30", b.BeforeTarget)

	// Upsert refreshes in place.
	require.NoError(t, s.UpsertWatchBaseline(ctx, "centralbarber", "34098477", "2026-07-20T09:00:00Z", ""))
	b, _, err = s.GetWatchBaseline(ctx, "centralbarber", "34098477")
	require.NoError(t, err)
	assert.Equal(t, "2026-07-20T09:00:00Z", b.NextAvailable)
}

func TestListBusinessesAndGet(t *testing.T) {
	ctx := context.Background()
	s := newVagaroTestStore(t)

	require.NoError(t, s.UpsertBusiness(ctx, BusinessRecord{
		Slug: "centralbarber", BusinessID: "93458", Name: "Central Barber",
		Rating: 4.9, ReviewCount: 212, PriceRange: "$$", City: "Seattle", State: "WA",
	}))
	b, ok, err := s.GetBusinessBySlug(ctx, "centralbarber")
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, 4.9, b.Rating)
	assert.Equal(t, "$$", b.PriceRange)

	_, ok, err = s.GetBusinessBySlug(ctx, "missing")
	require.NoError(t, err)
	assert.False(t, ok)

	all, err := s.ListBusinesses(ctx)
	require.NoError(t, err)
	require.Len(t, all, 1)
	assert.Equal(t, "Central Barber", all[0].Name)
}
