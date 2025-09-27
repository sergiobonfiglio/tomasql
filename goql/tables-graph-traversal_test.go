package goql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDBGraphData_ShortestJoinPaths(t *testing.T) {
	t.Run("joins(A, {C}) for graph A -> B -> C should be correct", func(t *testing.T) {
		// Setup mock graph: A -> B -> C
		A := mockTable("A")
		B := mockTable("B")
		C := mockTable("C")
		relationships := map[Table]map[Table]map[Column]Column{
			A: {
				B: {NewCol[int]("colAB", A): NewCol[int]("colBA", B)},
			},
			B: {
				C: {NewCol[int]("colBC", B): NewCol[int]("colCB", C)},
			},
		}

		g := &DBGraphData{relationships: relationships}

		joinPath, err := g.MinimalJoins(A, []Table{C})
		require.NoError(t, err)
		require.Len(t, joinPath, 2)

		require.Equal(t, B, joinPath[0].Target)
		require.Equal(t, C, joinPath[1].Target)
	})

	t.Run("joins(A, {C,D}) for graph A -> B -> C,D should be correct", func(t *testing.T) {
		// Setup mock graph: A -> B -> C,D
		A := mockTable("A")
		B := mockTable("B")
		C := mockTable("C")
		D := mockTable("D")
		relationships := map[Table]map[Table]map[Column]Column{
			A: {
				B: {NewCol[int]("colAB", A): NewCol[int]("colBA", B)},
			},
			B: {
				C: {NewCol[int]("colBC", B): NewCol[int]("colCB", C)},
				D: {NewCol[int]("colBD", B): NewCol[int]("colDB", D)},
			},
		}

		g := &DBGraphData{relationships: relationships}

		joinPath, err := g.MinimalJoins(A, []Table{C, D})
		require.NoError(t, err)

		require.Equal(t, B, joinPath[0].Target)
		require.Equal(t, C, joinPath[1].Target)
		require.Equal(t, D, joinPath[2].Target)
	})

	t.Run("joins(A, {C}) for graph A -> B -> C.col1, C.col2 should fail", func(t *testing.T) {
		// Setup mock graph: A -> B -> C
		A := mockTable("A")
		B := mockTable("B")
		C := mockTable("C")
		relationships := map[Table]map[Table]map[Column]Column{
			A: {
				B: {NewCol[int]("colAB", A): NewCol[int]("colBA", B)},
			},
			B: {
				C: {
					NewCol[int]("colBC", B):  NewCol[int]("colCB", C),
					NewCol[int]("colBC2", B): NewCol[int]("colCB2", C),
				},
			},
		}

		g := &DBGraphData{relationships: relationships}

		_, err := g.MinimalJoins(A, []Table{C})
		require.Error(t, err)
	})

	t.Run("joins(A, {B}) for graph A -> B -> C.col1, C.col2 should be correct", func(t *testing.T) {
		// Setup mock graph: A -> B -> C
		A := mockTable("A")
		B := mockTable("B")
		C := mockTable("C")
		relationships := map[Table]map[Table]map[Column]Column{
			A: {
				B: {NewCol[int]("colAB", A): NewCol[int]("colBA", B)},
			},
			B: {
				C: {
					NewCol[int]("colBC", B):  NewCol[int]("colCB", C),
					NewCol[int]("colBC2", B): NewCol[int]("colCB2", C),
				},
			},
		}

		g := &DBGraphData{relationships: relationships}

		joinPath, err := g.MinimalJoins(A, []Table{B})
		require.NoError(t, err)
		require.Len(t, joinPath, 1)
		require.Equal(t, B, joinPath[0].Target)
	})

	t.Run("joins(A, {C}) for graph A -> B | C should fail", func(t *testing.T) {
		// Setup mock graph: A -> B -> C
		A := mockTable("A")
		B := mockTable("B")
		C := mockTable("C")
		relationships := map[Table]map[Table]map[Column]Column{
			A: {
				B: {NewCol[int]("colAB", A): NewCol[int]("colBA", B)},
			},
			B: {
				A: {
					NewCol[int]("colBA", A): NewCol[int]("colAB", B),
				},
			},
			C: {
				A: {
					NewCol[int]("colCA", C): NewCol[int]("colAC", A),
				},
			},
		}

		g := &DBGraphData{relationships: relationships}

		_, err := g.MinimalJoins(A, []Table{C})
		require.Error(t, err)
	})
}
