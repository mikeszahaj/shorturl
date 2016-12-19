# Short URL

Proof of concept short URL service written in Go.  Has basic analytics including timestamps and browsers per short URL.

Uses [Gorm](https://github.com/jinzhu/gorm) and MySQL.

To use, create a `shorturls` database in MySQL.  Insert records into the `redirects` table.  Go to `http://localhost:8080/1a` where `1a` is the base-36 representation of an `id` in the redirects table (1a = 46).
