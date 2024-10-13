# click: Golang ClickHouse query SQL builder

## 1. design overview

- **declarative builder**: build ClickHouse SQL declaratively, in simple or advanced ways
- **type-safe**: static-typed, extensible SQL template for code reuse
- **feel-like-home**: SQL-like, idiomatic Go, no bloated stuff
- **easy to extend & interoperate with**:
    + good interoperability with existing popular ORM library (`github.com/huandu/go-sqlbuilder`)
    + directly generate final query SQL, easy to use outside the library

## 2. features & modules

1. `querybuilder`:
    + `SimpleQuery`: non-nested query shortcut, build with struct
    + `Select()`: declarative, chained, freestyle builder
2. `expression`: fundamental SQL expressions and operators
    + expressions:
    + operators:
3. `datasource`: datasource connector abstraction
    + `FromGoSQL()`: adaptor for ClickHouse connection from `sql` package
    + `Datasource`: interface where user can implement to use their own ClickHouse gateways, proixes...

## 2. examples

ˇ

### 2.1 click 101

```go
package main

import (
	"fmt"
	"github.com/keuin/click"
)

const (
	Date click.Column = "date"
	User click.Column = "user"
)

func main() {
	sumVisitCount := click.Alias("visit_count")
	sql, _ := click.
		Select(Date, User, click.As(click.Count(), sumVisitCount)).
		From("user_accesses").
		GroupBy(Date, User).
		OrderBy(sumVisitCount).
		Limit(10).PrettyPrint().BuildString()
	fmt.Println(sql)
}
```

Prints:

```clickhouse
SELECT date,
       user,
       count() AS visit_count
FROM user_accesses
GROUP BY date,
         user
ORDER BY visit_count
LIMIT 10
```

### 2.2 nested query

(TODO)
