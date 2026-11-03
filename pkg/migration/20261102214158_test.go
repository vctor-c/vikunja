// Vikunja is a to-do list application to facilitate your life.
// Copyright 2018-present Vikunja and contributors. All rights reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package migration

import (
	"testing"

	"code.vikunja.io/api/pkg/db"
	"github.com/stretchr/testify/require"
)

type projectViewBefore20261102214158 struct {
	ID           int64  `xorm:"autoincr not null unique pk"`
	ProjectID    int64  `xorm:"not null index"`
	Title        string `xorm:"varchar(255) not null"`
	DoneBucketID int64  `xorm:"bigint INDEX null"`
}

func (projectViewBefore20261102214158) TableName() string { return "project_views" }

func TestInProgressBucketMigration20261102214158(t *testing.T) {
	x, err := db.CreateTestEngine()
	require.NoError(t, err)
	table := projectViewBefore20261102214158{}
	require.NoError(t, x.DropTables(table))
	t.Cleanup(func() { require.NoError(t, x.DropTables(table)) })
	require.NoError(t, x.Sync2(table))
	_, err = x.Insert(&projectViewBefore20261102214158{
		ID:           1,
		ProjectID:    1,
		Title:        "Existing board",
		DoneBucketID: 3,
	})
	require.NoError(t, err)
	before, err := x.DBMetas()
	require.NoError(t, err)
	for _, m := range migrations {
		if m.ID == "20261102214158" {
			require.NoError(t, m.Migrate(x))
		}
	}
	after, err := x.DBMetas()
	require.NoError(t, err)
	found := false
	for _, oldTable := range before {
		if oldTable.Name != "project_views" {
			continue
		}
		for _, newTable := range after {
			if newTable.Name != oldTable.Name {
				continue
			}
			found = true
			require.NotNil(t, newTable.GetColumn("in_progress_bucket_id"))
			for _, column := range oldTable.ColumnsSeq() {
				require.NotNil(t, newTable.GetColumn(column))
			}
			for name, index := range oldTable.Indexes {
				preserved, exists := newTable.Indexes[name]
				require.True(t, exists)
				require.Equal(t, index.Type, preserved.Type)
				require.Equal(t, index.Cols, preserved.Cols)
			}
		}
	}
	require.True(t, found)
	stored := &projectViewBefore20261102214158{}
	exists, err := x.ID(1).Get(stored)
	require.NoError(t, err)
	require.True(t, exists)
	require.Equal(t, "Existing board", stored.Title)
	require.Equal(t, int64(3), stored.DoneBucketID)
	added := &projectViews20261102214158{}
	exists, err = x.Table("project_views").Where("id = ?", 1).Get(added)
	require.NoError(t, err)
	require.True(t, exists)
	require.Zero(t, added.InProgressBucketID)
}
