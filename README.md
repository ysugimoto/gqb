# gqb - Golang Simple Query Builder

## Features
- Build SQL easily through the method chains
- Returns abstact scanned result
- Query results can marshal JSON directly

## Installation

```shell
go get -u github.com/ysugimoto/gqb
```

## Usage

Example database is here (MySQL):

```sql
CREATE TABLE IF NOT EXISTS companies (
  id int(11) unsigned NOT NULL AUTO_INCREMENT,
  name varchar(255) NOT NULL,
  PRIMARY KEY (id)
) DEFAULT CHARSET=utf8;

INSERT INTO companies (name) VALUES ('Google'), ('Apple'), ('Microsoft');

CREATE TABLE IF NOT EXISTS company_attributes (
  id int(11) unsigned NOT NULL AUTO_INCREMENT,
  company_id int(11) unsigned NOT NULL,
  url varchar(255) NOT NULL,
  PRIMARY KEY (id)
) DEFAULT CHARSET=utf8;

INSERT INTO company_attributes (company_id, url) VALUES (1, 'https://google.com'), (2, 'https://apple.com'), (3, 'https://microsoft.com');
```

And make sure `*sql.DB` is created properly:

```go
// connect database as you expected
db, err := sql.Open("mysql", "user:pass@tcp(127.0.0.1:3306)/db_name")
if err != nil {
  log.Fatal(err)
}
defer db.Close()
```

### Getting started

The following example maybe generic usage. We expects SQL as:

```sql
SELECT name FROM companies WHERE id = 3;
```

`gqb` makes above SQL and retrieve result by following code:

```go
results, err := gqb.New(db).
  Select("name").
  Where("id", 1, gqb.Equal).
  Get("companies")

if err != nil {
  log.Fatal(err)
}
// retrieve result
for _, r := range results {
  fmt.Println(r.MustString("name")) //=> Google
}

// Also can marshal JSON directly
buf, _ := json.Marshal(results)
fmt.Println(string(buf)) //=> [{"name":"Google"}]

// Map to your struct
type Company struct {
  Name string `db:"name"`  // gqb maps value corresponds to "db" tag field
}
companies := []Company{}
if err := results.Map(&companies); err != nil {
  log.Fatal(err)
}
fmt.Println(companies[0].Name) //=> Google
```

If you want to get a single record, you can call `GetOne("companies")` instead.
To leaan more example usage, see [examples](https://github.com/ysugimoto/gqb/tree/master/examples).

## Query Execution

Note that `gqb` is just only for query bulder, so query exection, prepared statement, escaping bind parameters depend on `databae/sql`.

`gqb.New(db)` of first argument accepts `gqb.Executor` interface which has a couple of methods:

- `QueryContext(ctx context.Context, query string, binds ...interface{})`
- `ExecContext(ctx context.Context, query string, binds ...interface{})`

It means you can use as same syntax in transaction. `gqb.new(*sql.Tx)` also valid.

## Scan value

The `gqb.Result` struct can access through the `XXX(column)` or `MustXXX(column)`.
For example, to retrieve `id int(11)` column, you should call `result.MustInt64("id")`.

Occasionally there is a case that result value `null`, then you can call `v, err := result.Int64("id")`.
The `err` is returned if column value doesnt' exist or `null`.

Also, you can confirm field value is `null` via `result.Nil("id")`. It returns `true` is value is `null`.

And, if you want to use query result as your specific struct, you can call `result.Map(&strcut)`.
it will map values to field which corresponds to tag value of `db:"field"`.

`gqb` supports following struct field types:

- string / \*string
- int / \*int
- int8 / \*int8
- int16 / \*int16
- int32 / \*int32
- int64 / \*int64
- uint / \*uint
- uint8 / \*uint8
- uint16 / \*uint16
- uint32 / \*uint32
- uint64 / \*uint64
- float32 / \*float32
- float64 / \*float64

`[]byte`, corresponds to `blob` type column not supported.yet.

## Benchmarks

Native SQL vs `gqb` Query Builder.

100 records:

```
BenchmarkNativeSQL-8                2000            696598 ns/op            1072 B/op         34 allocs/op
BenchmarkQueryBuilder-8             2000            653312 ns/op            2910 B/op         87 allocs/op
```

1000 records:

```
BenchmarkNativeSQL-8                2000            738930 ns/op            1076 B/op         34 allocs/op
BenchmarkQueryBuilder-8             2000            681146 ns/op            2912 B/op         87 allocs/op
```

10000 records:

```
BenchmarkNativeSQL-8                2000            747242 ns/op            1073 B/op         34 allocs/op
BenchmarkQueryBuilder-8             2000            751494 ns/op            2914 B/op         87 allocs/op
```

## Author

Yoshiaki Sugimoto

## License

MIT

