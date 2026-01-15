package tomasql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddJoins(t *testing.T) {
	t.Run("empty join items", func(t *testing.T) {
		// Test with no join items - should return the same builder
		acc := Account.As("a")
		builder := Select(acc.Id).From(acc)
		result := _addJoins(builder, []*JoinItem{}...)

		sql, _ := result.SQL()
		expected := "SELECT a.id FROM account AS a"
		require.Equal(t, expected, sql)
	})

	t.Run("single join item", func(t *testing.T) {
		// Test with a single join item
		acc := Account.As("a")
		builder := Select(acc.Id, ShoppingCart.Id).From(acc)

		joinItem := &JoinItem{
			Target:      ShoppingCart,
			OnCondition: acc.Id.Eq(ShoppingCart.OwnerId),
		}

		result := _addJoins(builder, joinItem)
		sql, _ := result.SQL()
		expected := "SELECT a.id, shopping_cart.id FROM account AS a JOIN shopping_cart ON a.id = shopping_cart.owner_id"
		require.Equal(t, expected, sql)
	})

	t.Run("multiple join items", func(t *testing.T) {
		// Test with multiple join items
		acc := Account.As("a")
		sc := ShoppingCart.As("sc")
		builder := Select(acc.Id, sc.Id).From(acc)

		joinItem1 := &JoinItem{
			Target:      sc,
			OnCondition: acc.Id.Eq(sc.OwnerId),
		}
		joinItem2 := &JoinItem{
			Target:      ShoppingCart.As("sc2"),
			OnCondition: acc.Id.Eq(ShoppingCart.As("sc2").OwnerId),
		}

		result := _addJoins(builder, joinItem1, joinItem2)
		sql, _ := result.SQL()
		expected := "SELECT a.id, sc.id FROM account AS a " +
			"JOIN shopping_cart AS sc ON a.id = sc.owner_id " +
			"JOIN shopping_cart AS sc2 ON a.id = sc2.owner_id"
		require.Equal(t, expected, sql)
	})

	t.Run("nil join items should be skipped", func(t *testing.T) {
		// Test that nil join items are skipped
		acc := Account.As("a")
		builder := Select(acc.Id, ShoppingCart.Id).From(acc)

		joinItem := &JoinItem{
			Target:      ShoppingCart,
			OnCondition: acc.Id.Eq(ShoppingCart.OwnerId),
		}

		result := _addJoins(builder, nil, joinItem, nil)
		sql, _ := result.SQL()
		expected := "SELECT a.id, shopping_cart.id FROM account AS a JOIN shopping_cart ON a.id = shopping_cart.owner_id"
		require.Equal(t, expected, sql)
	})

	t.Run("Joins method uses _addJoins", func(t *testing.T) {
		// Test that the Joins method correctly delegates to _addJoins
		acc := Account.As("a")
		sc := ShoppingCart.As("sc")

		joinItem1 := &JoinItem{
			Target:      sc,
			OnCondition: acc.Id.Eq(sc.OwnerId),
		}

		sql, _ := Select(acc.Id, sc.Id).
			From(acc).
			Joins(joinItem1).
			SQL()

		expected := "SELECT a.id, sc.id FROM account AS a JOIN shopping_cart AS sc ON a.id = sc.owner_id"
		require.Equal(t, expected, sql)
	})

	t.Run("Joins with nil condition", func(t *testing.T) {
		// Test join item without condition
		acc := Account.As("a")
		builder := Select(acc.Id, ShoppingCart.Id).From(acc)

		joinItem := &JoinItem{
			Target:      ShoppingCart,
			OnCondition: nil,
		}

		result := _addJoins(builder, joinItem)
		sql, _ := result.SQL()
		expected := "SELECT a.id, shopping_cart.id FROM account AS a JOIN shopping_cart"
		require.Equal(t, expected, sql)
	})
}
