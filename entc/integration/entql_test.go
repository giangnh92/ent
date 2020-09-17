// Copyright 2019-present Facebook Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package integration

import (
	"context"
	"testing"

	"github.com/facebook/ent/entc/integration/ent"
	"github.com/facebook/ent/entc/integration/ent/pet"
	"github.com/facebook/ent/entql"

	"github.com/stretchr/testify/require"
)

func EntQL(t *testing.T, client *ent.Client) {
	require := require.New(t)
	ctx := context.Background()

	a8m := client.User.Create().SetName("a8m").SetAge(30).SaveX(ctx)
	nati := client.User.Create().SetName("nati").SetAge(30).AddFriends(a8m).SaveX(ctx)

	uq := client.User.Query()
	uq.EntQL().Where(entql.HasEdge("friends"))
	require.Equal(2, uq.CountX(ctx))

	uq = client.User.Query()
	uq.EntQL().Where(
		entql.And(
			entql.FieldEQ("name", "nati"),
			entql.HasEdge("friends"),
		),
	)
	require.Equal(nati.ID, uq.OnlyIDX(ctx))

	xabi := client.Pet.Create().SetName("xabi").SetOwner(a8m).SaveX(ctx)
	luna := client.Pet.Create().SetName("luna").SetOwner(nati).SaveX(ctx)
	uq = client.User.Query()
	uq.EntQL().Where(
		entql.And(
			entql.HasEdge("pets"),
			entql.HasEdgeWith("friends", entql.FieldEQ("name", "nati")),
		),
	)
	require.Equal(a8m.ID, uq.OnlyIDX(ctx))
	uq = client.User.Query()
	uq.EntQL().Where(
		entql.And(
			entql.HasEdgeWith("pets", entql.FieldEQ("name", "luna")),
			entql.HasEdge("friends"),
		),
	)
	require.Equal(nati.ID, uq.OnlyIDX(ctx))

	uq = client.User.Query()
	uq.EntQL().WhereName(entql.StringEQ("a8m"))
	require.Equal(a8m.ID, uq.OnlyIDX(ctx))
	pq := client.Pet.Query()
	pq.EntQL().WhereName(entql.StringOr(entql.StringEQ("xabi"), entql.StringEQ("luna")))
	require.Equal([]int{luna.ID, xabi.ID}, pq.Order(ent.Asc(pet.FieldName)).IDsX(ctx))

	pq = client.Pet.Query()
	pq.Where(pet.Name(luna.Name)).EntQL().WhereID(entql.IntEQ(luna.ID))
	require.Equal(luna.ID, pq.Order(ent.Asc(pet.FieldName)).OnlyIDX(ctx))
	pq = client.Pet.Query()
	pq.Where(pet.Name(luna.Name)).EntQL().WhereID(entql.IntEQ(xabi.ID))
	require.False(pq.ExistX(ctx))
}
