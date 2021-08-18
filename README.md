# go-sqljs-driver

## Caveats

### Named parameters prefixes

In [SQLite's named parameters](https://www.sqlite.org/lang_expr.html#varparam) the prefix (`:`, `@`, or `$`) is included as part of the name.

However [Go's `database/sql.NamedArg` type](https://pkg.go.dev/database/sql#NamedArg) specifies that *"Name must omit any symbol prefix."*.

This makes it impossible for `go-sqljs-driver` to bind a different value to named parameters with the same name but a different prefix:

```sql
SELECT * FROM example WHERE col1 = :param1 AND col2 = @param1
-- :param1 and @param1 will always have the same value
```
