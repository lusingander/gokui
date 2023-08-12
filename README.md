# gokui

my sql utilities

## Installation

```
$ go install github.com/lusingander/gokui@latest
```

(require Go 1.20+)

## Usage

### `generate` command

```
$ cat create.sql 
CREATE TABLE users (
    user_id VARCHAR(256) NOT NULL COMMENT 'user id',
    name VARCHAR(256),
    age INT UNSIGNED NOT NULL,
    created_at DATETIME(6) NOT NULL,
    primary key (user_id)
);

$ cat create.sql | gokui generate select --newline
SELECT
  user_id,
  name,
  age,
  created_at
FROM
  users
;

$ cat create.sql | gokui generate insert --newline --insert-select
INSERT INTO users
(
  user_id,
  name,
  age,
  created_at
)
SELECT
  user_id,
  name,
  age,
  created_at
FROM
  users
;
```

## License

MIT
