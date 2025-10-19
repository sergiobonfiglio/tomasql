package postgres

import (
	"fmt"

	"github.com/sergiobonfiglio/tomasql"
	"github.com/sergiobonfiglio/tomasql/dialects/pgres"
)

//go:generate go run github.com/sergiobonfiglio/tomasql/cmd/table-def-gen --schema ../schema.sql --package-dir . --package-name postgres --table-graph-file= --with-pgres-extensions

func Run() {

	// Set PostgreSQL dialect
	pgres.SetDialect()

	fmt.Printf("DIALECT: %s",
		tomasql.GetDialect().Name())

	// Example 1: Simple SELECT query
	fmt.Println("\n--- Example 1: Basic SELECT ---")
	sql, params := tomasql.Select(Users.Id, Users.Name, Users.Email).
		SQL()
	fmt.Printf("SQL: %s\n", sql)
	fmt.Printf("Params: %v\n", params)

	// Example 2: SELECT with FROM
	fmt.Println("\n--- Example 2: SELECT with FROM ---")
	sql, params = tomasql.Select(Users.Id, Users.Name, Users.Email).
		From(Users).
		SQL()
	fmt.Printf("SQL: %s\n", sql)
	fmt.Printf("Params: %v\n", params)

	// Example 3: SELECT with WHERE
	fmt.Println("\n--- Example 3: SELECT with WHERE ---")
	sql, params = tomasql.Select(Users.Id, Users.Name, Users.Email).
		From(Users).
		Where(Users.IsActive.EqParam(true)).
		SQL()
	fmt.Printf("SQL: %s\n", sql)
	fmt.Printf("Params: %v\n", params)

	// Example 4: JOIN query
	fmt.Println("\n--- Example 4: JOIN query ---")
	sql, params = tomasql.Select(Users.Name, Orders.Id.As("order_id"), Orders.TotalAmount).
		From(Users).
		Join(Orders).On(Users.Id.Eq(Orders.UserId)).
		Where(Users.IsActive.EqParam(true)).
		SQL()
	fmt.Printf("SQL: %s\n", sql)
	fmt.Printf("Params: %v\n", params)

	// Example 5: Complex query with aliases and multiple conditions
	fmt.Println("\n--- Example 5: Complex query ---")
	u := Users.As("u")
	o := Orders.As("o")
	p := Products.As("p")

	sql, params = tomasql.Select(u.Name.As("customer_name"), o.Id.As("order_id"), p.Name.As("product_name")).
		From(u).
		Join(o).On(u.Id.Eq(o.UserId)).
		Join(p).On(p.Id.EqParam(1)). // Simplified join for example
		Where(u.IsActive.EqParam(true).And(o.Status.EqParam("completed"))).
		OrderBy(u.Name.Asc(), o.Id.Desc()).
		SQL()
	fmt.Printf("SQL: %s\n", sql)
	fmt.Printf("Params: %v\n", params)

	// Example 6: Simple aggregation query
	fmt.Println("\n--- Example 6: Aggregation query ---")
	sql, params = tomasql.Select(Users.Name, tomasql.Count().As("order_count")).
		From(Users).
		LeftJoin(Orders).On(Users.Id.Eq(Orders.UserId)).
		Where(Users.IsActive.EqParam(true)).
		SQL()
	fmt.Printf("SQL: %s\n", sql)
	fmt.Printf("Params: %v\n", params)

	// Example 7: Using PostgreSQL specific column type
	fmt.Println("\n--- Example 7: PostgreSQL specific column type ---")
	sql, params = tomasql.Select(Users.CreatedAt, Users.UpdatedAt).
		From(Users).
		Where(Users.Name.ILikeParam("%test.com")).
		SQL()
	fmt.Printf("SQL: %s\n", sql)
	fmt.Printf("Params: %v\n", params)

}
