# gorm connects postgres behind socks5 proxy


```
❯ pg_ctl init -D ./data
❯ pg_ctl start -D ./data
❯ createuser -s postgres
❯ go run .

❯ psql -h localhost -U postgres
postgres=#
postgres=# select * from books;
 id |                        title                        |  author
----+-----------------------------------------------------+----------
  1 | 2024-12-03 17:17:40.984556 +0800 CST m=+0.040438751 | fixpoint
(1 row)

postgres=#

```
