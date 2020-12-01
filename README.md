# SQL WHERE clause builder for Go

`gowhere` is neither an ORM package nor full-featured SQL builder. It only provides a flexible and powerful way to build SQL WHERE string from a simple array/slice or a map.

The goal of this package is to create an adapter for frontend app to easily "query" the given data. In the other way around, backend app can "understand" the request from frontend app to correctly build the SQL query. With minimum configurations and coding!

This package can also be used together with standard "database/sql" package or other ORM packages for better experience, if any :)

## Install

```bash
go get -u github.com/imdatngo/gowhere
```

This package is dependency-free!

## Usages

The simple way:

```go
import "github.com/imdatngo/gowhere"

plan := gowhere.Where(map[string]interface{}{
    "name__contains": "Gopher",
    "budget__gte": 1000,
    "date__between": []interface{}{"2019-04-13", "2019-04-15"},
})

plan.SQL()
// ("name" LIKE ? AND "budget" >= ? and "date" BETWEEN ? AND ?)

plan.Vars()
// [%Gopher% 1000 2019-04-13 2019-04-15]
```

The advanced way:

```go
import "github.com/imdatngo/gowhere"

// initialize with all available configurations
plan := gowhere.WithConfig(gowhere.Config{
    Separator: "__",
    Dialect: gowhere.DialectPostgreSQL,
    Strict: true,
    Table: "",
    ColumnAliases: map[string]string{},
    CustomConditions: map[string]CustomConditionFn{
        "search": func(key string, val interface{}, cfg *gowhere.Config) interface{} {
            val = "%" + val.(string) + "%"
            return []interface{}{"lower(full_name) like ? OR lower(title) like ?", val, val}
        },
    },
})

// modify the config on the demands
plan.SetTable("trips").SetColumnAliases(map[string]string{"name": "full_name"})

// a map will be translated to "AND" conditions
plan.Where(map[string]interface{}{
    "name__contains": "Gopher",
    "budget__gte": 1000,
})

// a slice will be "OR" conditions
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

## Operator

For example: `"name__startswith"`, `name` is the field(column) and `startswith` is the operator. Django developer might find this familiar ;)

Operator is an user friendly name for a specific SQL operator. It's a suffix which added into the field to reduce the complexity from input schema, yet flexible enough to generate complex conditions.

Operator is optional. If not given, it's set to:

- `isnull` if the value is `nil`. E.g: `{"name": nil}` => sql, vars: `name IS NULL, []`
- `in` if the value is slice or array. E.g: `{"id": []int{1, 2, 3}}` => `id in (?), [[1 2 3]]`
- `exact` if otherwise. E.g: `{"name": "Gopher"}` => `name = ?, [Gopher]`

Built-in operators:

- `exact`: Exact match, using `=` operator.
- `iexact`: Case-insensitive exact match, wrap both column and value with `lower()` function.
- `notexact`: Opposite of `exact`
- `notiexact`: Opposite of `iexact`
- `gt`: Greater than
- `gte`: Greater than or equal to
- `lt`: Less than
- `lte`: Less than or equal to
- `startswith`: Case-sensitive starts-with, auto cast value to string with `%` suffix
- `istartswith`: Case-insensitive starts-with
- `endswith`: Case-sensitive ends-with, auto cast value to string with `%` prefix
- `iendswith`: Case-insensitive ends-with
- `contains`: Case-insensitive containment test, auto cast value to string with both `%` suffix, prefix
- `icontains`: Case-sensitive containment test
- `in`: In a given slice, array
- `date`: For datetime fields, casts the value as date
- `between`: For datetime string fields, range test
- `isnull`: Takes either True or False, which correspond to SQL
  queries of IS NULL and IS NOT NULL, respectively.
- `datebetween`: For query datetime range fields

## TODO

- [x] Publish!
- [x] Ability to add custom operators
- [x] Ability to add custom conditions
- [ ] Full tests with 100% code coverage
- [ ] Manipulate the conditions? Such as `HasCondition()`, `UpdateCondition()`, `RemoveCondition()`?

## License

© Dat Ngo, 2019~time.Now()

Released under the [MIT License](./LICENSE)
