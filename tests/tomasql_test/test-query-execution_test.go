package tomasql_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	. "github.com/sergiobonfiglio/tomasql"
	"github.com/sergiobonfiglio/tomasql/cmd/table-def-gen/setup"
	"github.com/sergiobonfiglio/tomasql/dialects/pgres"
	"github.com/stretchr/testify/require"
)

func TestQueryOnDB(t *testing.T) {

	db, err := setup.SetupTestContainer(t,
		"example_schema.sql",
		"postgres:latest")
	if err != nil {
		t.Fatalf("failed to set up test container: %v", err)
	}
	defer func() { _ = db.Close() }()

	type test struct {
		got  string
		want string
	}

	pgres.SetDialect()
	t.Run("subqueries with params", func(t *testing.T) {
		testItem := test{
			want: "SELECT * FROM (SELECT config.uuid FROM config WHERE config.uuid = $1) AS t1 " +
				"JOIN (SELECT config.uuid FROM config WHERE config.uuid = $2) AS t2 ON t1.uuid = t2.uuid " +
				"JOIN (SELECT config.uuid FROM config WHERE config.uuid = $2) AS t3 ON t2.uuid = t3.uuid",
			got: func() string {
				t1 := Select(Config.Uuid).From(Config).Where(Config.Uuid.EqParam("1")).AsNamedSubQuery("t1")
				t2 := Select(Config.Uuid).From(Config).Where(Config.Uuid.EqParam("2")).AsNamedSubQuery("t2")
				t3 := Select(Config.Uuid).From(Config).Where(Config.Uuid.EqParam("2")).AsNamedSubQuery("t3")

				t1Col := NewCol[string]("uuid", t1)
				t2Col := NewCol[string]("uuid", t2)
				t3Col := NewCol[string]("uuid", t3)

				sql, params := SelectAll().
					From(t1).
					Join(t2).On(t1Col.Eq(t2Col)).
					Join(t3).On(t2Col.Eq(t3Col)).
					SQL()

				require.Len(t, params, 2) // only 2 distinct params
				return sql
			}(),
		}

		require.Equal(t, testItem.want, testItem.got)
	})

	// setup
	type accountRow struct {
		Id        int64  `db:"id"`
		Uuid      string `db:"uuid"`
		CreatedTs int64  `db:"created_ts"`
	}
	var inserted []accountRow
	for range 10 {

		accUuid := uuid.NewString()
		createdTs := time.Now().Unix()
		accId := new(int64)
		err := db.Get(accId,
			`INSERT INTO account (uuid, created_ts) VALUES ($1, $2) RETURNING id`,
			accUuid,
			createdTs,
		)
		require.NoError(t, err)
		require.NotNil(t, accId)

		inserted = append(inserted, accountRow{
			Id:        *accId,
			Uuid:      accUuid,
			CreatedTs: createdTs,
		})
	}

	t.Run("query with params", func(t *testing.T) {
		builder := Select(Account.Id, Account.Uuid, Account.CreatedTs).
			From(Account).
			Where(Account.Id.EqParam(inserted[0].Id).
				And(Account.Id.EqParam(inserted[0].Id)).
				And(Account.Uuid.EqParam(inserted[0].Uuid)))

		res := accountRow{}
		sql, params := builder.SQL()
		require.Len(t, params, 2) // only 2 distinct params

		err := db.Get(&res, sql, params...)
		require.NoError(t, err)

		require.Equal(t, inserted[0].Id, res.Id)
		require.Equal(t, inserted[0].Uuid, res.Uuid)
		require.Equal(t, inserted[0].CreatedTs, res.CreatedTs)
	})

}
