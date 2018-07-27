`gqb` is a simple Query Builder for Golang.

- Build SQL easily through the method chains
- `gqb` retunrs abstact scanned result. Especially it's useful for _JOIN-ed_ query result
- Query results can marshal JSON directly

## Installation

```
go get -u github.com/ysugimoto/gqb
```

## Support database

Now we only tested on MySQL.

## Usage

Example database is here (MySQL):

```
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

INSERT INTO companie_attributes (company_id, url) VALUES (1, 'https://google.com'), (2, 'https://apple.com'), (3, 'https://microsoft.com');
```

And make sure `*sql.DB` is created propery:

```
// connect database as you expected
db, err := sql.Open("mysql", "user:pass@tcp(127.0.0.1:3306)/db_name")
if err != nil {
  log.Fatal(err)
}
defer db.Close()
```

### Select query

#### General SELECT query usage

```
results, err := gqb.New(db).
  Select("name").
  Where("id", 3, gqb.Lt).
  OrderBy("id", gqb.Desc).
  Get("companies")

if err != nil {
  log.Fatal(err)
}
// retrieve result
for r := range results {
  fmt.Println(r.MustString("name")) //=> Apple on first, Goolge on second
}

// Also can marshal JSON directly
buf, _ := json.Marshal(results)
fmt.Printf(string(buf)) //=> [{"name":"Apple"},{"name":"Google"}]

// Map to your struct
type Company struct {
  Name string `db:"name"`  // gqb maps value corresponds to "db" tag field
}
companies := []Company{}
if err := result.Map(&companies); err != nil {
  log.Fatal(err)
}
fmt.Println(companies[0].Name) //=> Apple
```

If you want to get a single record, you can call `GetOne("companies")` instead.

#### Use JOIN case

```
result, err := gqb.New(db).
  Select("name", "url").
  Join("company_attributes", "company_id", "id", gqb.Equal).
  Where("id", 3, gqb.Equal)
  GetOne("companies")

if err != nil {
  log.Fatal(err)
}
// retrieve result
fmt.Println(result.MustString("url")) //=> https://microsoft.com

// Also can marshal JSON directly
buf, _ := json.Marshal(result)
fmt.Printf(string(buf)) //=> {"name":"Microsoft","url":"https://microsoft.com"}

// Map to your struct
type Company struct {
  Name string `db:"name"`
  Url string `db:"url"`
}
ms := Company{}
if err := result.Map(&ms); err != nil {
  log.Fatal(err)
}
fmt.Println(ms.Name) //=> Microsoft
```

#### Insert / Update

for `INSERT` or `UPDATE`, `gqb` supplied `gqb.Data` which is shorthand alias for `map[string]interface{}`:

```
data := gqb.Data{
  "name": "Slakc", // will be update to correct name :-)
}
// the result is sql.Result
result, err := gqb.New(db).Insert("companies", data)

if err != nil {
  log.Fatal(err)
}

// Oops, update company name to correct one!
id, err := result.LastInsertId()
if err != nil {
  log.Println("failed to retrieve last inserted ID")
  return
}

data = gqb.Data{
  "name: "Slack",
}
result, err := gqb.New(db).
  Where("id", id, gqb.Eaual).
  Update("companies", data)

if err != nil {
  log.Fatal(err)
}
// Copany name will be updated to "Slack"!
```
