# SQL WHERE clause builder for Go #

`gowhere` is neither an ORM package nor full-featured SQL builder. It only provides a flexible and powerful way to build SQL WHERE string from a simple array/slice or a map.

The goal of this package is to create an adapter for frontend app to easily "query" the given data. In the other way around, backend app can "understand" the request from frontend app to correctly build the SQL query. With minimum configurations and coding!

This package can also be used together with standard "database/sql" package or other ORM packages for better experience, if any :)

## Install ##

```bash
go get -u github.com/imdatngo/gowhere
```

This package is dependency-free!

## Usages ##

The simple way:

```go
import "github.com/imdatngo/gowhere"

plan := gowhere.Where(map[string]interface{}{
    "name_contains": "Gopher",
    "budget__gte": 1000,
    "date__between": []interface{}{"2019-04-13", "2019-04-15"},
})

plan.SQL()
// ("name" = ? AND "budget" >= ? and "date" BETWEEN ? AND ?)

plan.Vars()
// [Gopher 1000 2019-04-13 2019-04-15]
```

The advanced way:

```go
import "github.com/imdatngo/gowhere"

// initialize with all available configurations
plan := gowhere.WithConfig(gowhere.Config{
    Separator: "__",
    Dialect: gowhere.DialectPostgreSQL,
    Strict: true,
    Table: "trips",
    ColumnAliases: map[string]string{"name": "full_name"},
})

// a map will be translated to "AND" conditions
plan.Where(map[string]interface{}{
    "name__contains": "Gopher",
    "budget__gte": 1000,
})

// a slice will be "OR" conditions, check here below for more advanced way to build this "OR" conditions
plan.Where([]map[string]interface{}{
    {"started_at__date": "2019-04-13"},
    {"started_at__date": "2019-04-15"},
    {"started_at__gte": time.Date(2019, 4, 19, 0, 0, 0, 0, time.Local)},
})

// same as `Where`, `Not` receives either map, slice or raw SQL string. Then it simply wraps the conditions with "NOT" keyword
plan.Not("members < ? AND members > ?", 2, 10)

// `Or` be like: ((all_current_conditions) OR (new_conditions))
plan.Or("anywhere = TRUE")

// In "Strict" mode, any invalid conditions/operators if given will cause `InvalidCond` error. Not to mention the not-tested runtime errors :)
if err := plan.Build().Error; err != nil {
    panic(err)
}

plan.SQL()
// ((("trips"."full_name" LIKE ? AND "trips"."budget" >= ?) AND ((DATE("trips"."started_at") = ?) OR (DATE("trips"."started_at") = ?) OR ("trips"."started_at" >= ?)) AND NOT (members < ? AND members > ?)) OR (anywhere = TRUE))

plan.Vars()
// [%Gopher% 1000 2019-04-13 2019-04-15 2019-04-19 2 10]
```

## TODO ##

- [x] Publish!
- [ ] Full tests with 100% code coverage
- [ ] Ability to add custom operators
- [ ] Manipulate the conditions? Such as `HasCondition()`, `UpdateCondition()`, `RemoveCondition()`?