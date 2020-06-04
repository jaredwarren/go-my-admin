package db

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"

	_ "github.com/go-sql-driver/mysql"
	sqlfmt "github.com/kanmu/go-sqlfmt/sqlfmt"
	"github.com/xwb1989/sqlparser"
)

// Query ...
type Query struct {
	stmt *sqlparser.Select
	args []interface{}
}

// NewQuery ...
func NewQuery(q string, args ...interface{}) (*Query, error) {
	// parse query
	s, err := sqlparser.Parse(q)
	if err != nil {
		return nil, err
	}
	var stmt *sqlparser.Select
	switch s.(type) {
	case *sqlparser.Select:
		stmt = s.(*sqlparser.Select)
	default:
		return nil, fmt.Errorf("unknown query type::%s", reflect.TypeOf(s))
	}
	return &Query{
		stmt: stmt,
		args: args,
	}, nil
}

// Sort add "Order By" to query
func (q *Query) Sort(col, dir string) {
	if dir != "" && col != "" {
		// get last OrderBy
		if len(q.stmt.OrderBy) > 0 {
			// update direction
			or := q.stmt.OrderBy[len(q.stmt.OrderBy)-1]
			or.Direction = dir

			// update name
			e := or.Expr.(*sqlparser.ColName)
			e.Name = sqlparser.NewColIdent(col)
		} else {
			// append OrderBy
			e := &sqlparser.ColName{}
			e.Name = sqlparser.NewColIdent(col)
			or := &sqlparser.Order{
				Direction: dir,
				Expr:      e,
			}
			q.stmt.OrderBy = append(q.stmt.OrderBy, or)
		}
	}
}

// String convert query statement to string
func (q *Query) String() string {
	return sqlparser.String(q.stmt)
}

// Search add filter (WHERE {searchcol} like {search}) to base query
func (q *Query) Search(search, searchcol string) {
	// Search
	if search != "" && searchcol != "" {
		if q.stmt.Where != nil {
			where := q.stmt.Where
			// Wrap old where in paren exp AND new like
			where.Expr = &sqlparser.AndExpr{
				Left: &sqlparser.ParenExpr{
					Expr: where.Expr,
				},
				Right: &sqlparser.ComparisonExpr{
					Operator: sqlparser.LikeStr,
					Left: &sqlparser.ColName{
						Name: sqlparser.NewColIdent(searchcol),
					},
					Right: &sqlparser.SQLVal{
						Type: 0,
						Val:  []byte(search),
					},
				},
			}
		} else {
			q.stmt.Where = &sqlparser.Where{
				Type: "where",
				Expr: &sqlparser.ComparisonExpr{
					Operator: sqlparser.LikeStr,
					Left: &sqlparser.ColName{
						Name: sqlparser.NewColIdent(searchcol),
					},
					Right: &sqlparser.SQLVal{
						Type: 0,
						Val:  []byte(search),
					},
				},
			}
		}
	}
}

// Result ...
type Result struct {
	Rows       []map[string]interface{}
	Columns    []string
	ColumnType []*sql.ColumnType
	Time       time.Duration
}

// Beautify clean up a query
func Beautify(query string) string {
	opt := &sqlfmt.Options{}
	clean, _ := sqlfmt.Format(query, opt)
	return clean
}
