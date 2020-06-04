package main

import (
	"fmt"

	"github.com/jaredwarren/app"
	"github.com/jaredwarren/myadmin/service"
	"github.com/xwb1989/sqlparser"

	_ "github.com/go-sql-driver/mysql"
)

func testStuff() {
	// s, err := sqlparser.Parse("select * from orders WHERE order_id like `asdf` or id=1")
	// s, err := sqlparser.Parse("select * from users where username like 'jwarr%'")
	// s, err := sqlparser.Parse("select * from orders WHERE order_id like `asdf`")
	s, err := sqlparser.Parse("select * from orders")
	if err != nil {
		panic(err)
	}
	stmt := s.(*sqlparser.Select)

	fmt.Printf("%+v\n", stmt)
	// fmt.Printf("%+v\n", stmt.Where.Expr)
	// fmt.Println(reflect.TypeOf(stmt.Where.Expr))

	// ex := stmt.Where.Expr.(*sqlparser.ComparisonExpr)
	// fmt.Printf("  %+v\n", ex.Right)
	// fmt.Println(reflect.TypeOf(ex.Right))

	// va := ex.Right.(*sqlparser.SQLVal)
	// fmt.Printf("     %s\n", va.Val)

	// return

	if stmt.Where != nil {
		where := stmt.Where
		// Wrap old where in paren exp AND new like
		where.Expr = &sqlparser.AndExpr{
			Left: &sqlparser.ParenExpr{
				Expr: where.Expr,
			},
			Right: &sqlparser.ComparisonExpr{
				Operator: sqlparser.LikeStr,
				Left: &sqlparser.ColName{
					Name: sqlparser.NewColIdent("user"),
				},
				Right: &sqlparser.ColName{
					Name: sqlparser.NewColIdent("jared"),
				},
			},
		}
	} else {
		stmt.Where = &sqlparser.Where{
			Type: "where",
			Expr: &sqlparser.ComparisonExpr{
				Operator: sqlparser.LikeStr,
				Left: &sqlparser.ColName{
					Name: sqlparser.NewColIdent("user"),
				},
				Right: &sqlparser.SQLVal{
					Type: 0,
					Val:  []byte("jarr%"),
				},
			},
		}
	}

	fmt.Println("\n\n:::", sqlparser.String(stmt))

	return

	// fmt.Println(reflect.TypeOf(where.Expr))
	// parnExp := where.Expr.(*sqlparser.ParenExpr)
	// fmt.Printf("\nP:%+v\n", parnExp)

	// fmt.Printf("\n %+v\n", parnExp.Expr)
	// fmt.Println(reflect.TypeOf(parnExp.Expr))

	// TODO: figure out how to wrap where.Expr into ParenExpr, then AND search

	// exp := where.Expr.(*sqlparser.AndExpr)
	// fmt.Printf("Expr:%+v\n", exp.Left)
	// fmt.Println(reflect.TypeOf(exp.Left))

	// ll := exp.Left.(*sqlparser.ComparisonExpr)
	// fmt.Printf("Expr L:%+v\n", ll.Left)
	// fmt.Println(reflect.TypeOf(ll.Right))
	// fmt.Printf("Expr R:%+v\n", ll.Right)
	// fmt.Println(reflect.TypeOf(ll.Right))

	// fmt.Printf("Expr:%+v\n", exp.Right)
	// fmt.Println(reflect.TypeOf(exp.Right))
	// fmt.Printf("Type:%+v\n", where.Type) //  = where

	// return

	// sl := &sqlparser.ColName{}
	// sl.Name = sqlparser.NewColIdent("user")

	// sr := &sqlparser.ColName{}
	// sr.Name = sqlparser.NewColIdent("jared")

	// right := &sqlparser.ComparisonExpr{
	// 	Operator: sqlparser.LikeStr,
	// 	Left:     sl,
	// 	Right:    sr,
	// }

	// andExp := &sqlparser.AndExpr{
	// 	Left:  where.Expr,
	// 	Right: right,
	// }

	// where.Expr = andExp

	// fmt.Println("\n\n:::", sqlparser.String(stmt))
	// return

	// // parse query
	// // stmt, err := sqlparser.Parse("select * from order_cats LIMIT 1, 2")
	// stmt, err := sqlparser.Parse("select * from order_cats")

	// if err != nil {
	// 	fmt.Println("E:", err)
	// 	return
	// }
	// // cast statement
	// st := stmt.(*sqlparser.Select)

	// // offset := sqlparser.NewIntVal([]byte("1"))
	// rowCount := sqlparser.NewIntVal([]byte("1"))
	// lim := &sqlparser.Limit{
	// 	// Offset:   offset,
	// 	Rowcount: rowCount,
	// }
	// st.Limit = lim

	// fmt.Println(reflect.TypeOf(st.Limit.Offset))

	// fmt.Printf("%+v\n", lim)
	// fmt.Printf("%+v\n", st.Limit)
	// fmt.Printf("%+v\n", st.Limit.Offset)
	// fmt.Printf("%+v\n", st.Limit.Rowcount)

	// fmt.Println(reflect.TypeOf(st.Limit.Offset))

	// get last OrderBy
	// e := &sqlparser.ColName{}
	// e.Name = sqlparser.NewColIdent("sortname")
	// or := &sqlparser.Order{
	// 	Direction: "ASC",
	// 	Expr:      e,
	// }
	// st.OrderBy = append(st.OrderBy, or)

	// or.Direction = "ASC"
	// fmt.Printf("%+v\n", or.Expr)
	// e := or.Expr.(*sqlparser.ColName)
	// e.Name = sqlparser.NewColIdent("ASDF")
	// exp := or.Expr.(*sqlparse)
	// or.Expr

	// for _, y := range st.OrderBy {
	// 	fmt.Printf("%+v\n", y.Expr)
	// }

	// fmt.Println(sqlparser.String(stmt))
	// fmt.Println("done")

	// // Otherwise do something with stmt
	// return

	// parse query
	// stmt, err := sqlparser.Parse("select * from order_cats ORDER BY name DESC")
	// if err != nil {
	// 	fmt.Println("E:", err)
	// 	return
	// }
	// // cast statement
	// st := stmt.(*sqlparser.Select)

	// // get last OrderBy
	// if len(st.OrderBy) > 0 {
	// 	or := st.OrderBy[len(st.OrderBy)-1]
	// 	or.Direction = "ASC"
	// 	fmt.Printf("%+v\n", or.Expr)
	// 	e := or.Expr.(*sqlparser.ColName)
	// 	e.Name = sqlparser.NewColIdent("ASDF")
	// 	// fmt.Println(reflect.TypeOf(e.Name))
	// 	// exp := or.Expr.(*sqlparse)
	// 	// or.Expr
	// }

	// // for _, y := range st.OrderBy {
	// // 	fmt.Printf("%+v\n", y.Expr)
	// // }

	// fmt.Println(sqlparser.String(stmt))
	// fmt.Println("done")

	// // Otherwise do something with stmt
	// return
}

func main() {
	// // load config
	// viper.SetConfigName("config_" + runtime.GOOS)
	// viper.AddConfigPath(".")
	// if err := viper.ReadInConfig(); err != nil {
	// 	log.Fatalf("Error reading config file, %s", err)
	// }

	// var serverConfig app.Config
	// viper.UnmarshalKey("server", &serverConfig)
	// fmt.Printf("%+v\n", serverConfig)
	// return

	conf := &app.WebConfig{
		Host: "127.0.0.1",
		Port: 8084,
	}
	a := app.NewWeb(conf)

	service.Register(a)

	d := <-a.Exit
	fmt.Printf("Done:%+v\n", d)
}
