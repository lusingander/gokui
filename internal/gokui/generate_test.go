package gokui

import (
	"fmt"
	"testing"
)

const (
	createTableUsers = `
CREATE TABLE users (
    user_id VARCHAR(256) NOT NULL COMMENT 'user id',
    name VARCHAR(256),
    age INT UNSIGNED NOT NULL,
    created_at DATETIME(6) NOT NULL,
    primary key (user_id)
);
`
)

func TestGenerateSelect(t *testing.T) {
	tests := []struct {
		sql  string
		opt  GenerateSelectOptions
		want string
	}{
		{
			sql: createTableUsers,
			opt: GenerateSelectOptions{
				NewLine: false,
			},
			want: `SELECT user_id, name, age, created_at FROM users;`,
		},
		{
			sql: createTableUsers,
			opt: GenerateSelectOptions{
				NewLine: true,
			},
			want: `SELECT
  user_id,
  name,
  age,
  created_at
FROM
  users
;`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%+v", tt.opt), func(t *testing.T) {
			got, err := GenerateSelect(tt.sql, tt.opt)
			if err != nil {
				t.Errorf("error occurred: %v", err)
			}
			if got != tt.want {
				t.Errorf("got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func TestGenerateInsert(t *testing.T) {
	tests := []struct {
		sql  string
		opt  GenerateInsertOptions
		want string
	}{
		{
			sql: createTableUsers,
			opt: GenerateInsertOptions{
				NewLine:      false,
				InsertSelect: false,
			},
			want: `INSERT INTO users (user_id, name, age, created_at) VALUES ('', '', 0, '');`,
		},
		{
			sql: createTableUsers,
			opt: GenerateInsertOptions{
				NewLine:      true,
				InsertSelect: false,
			},
			want: `INSERT INTO users
(
  user_id,
  name,
  age,
  created_at
)
VALUES
(
  '',
  '',
  0,
  ''
);`,
		},
		{
			sql: createTableUsers,
			opt: GenerateInsertOptions{
				NewLine:      false,
				InsertSelect: true,
			},
			want: `INSERT INTO users (user_id, name, age, created_at) SELECT user_id, name, age, created_at FROM users;`,
		},
		{
			sql: createTableUsers,
			opt: GenerateInsertOptions{
				NewLine:      true,
				InsertSelect: true,
			},
			want: `INSERT INTO users
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
;`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%+v", tt.opt), func(t *testing.T) {
			got, err := GenerateInsert(tt.sql, tt.opt)
			if err != nil {
				t.Errorf("error occurred: %v", err)
			}
			if got != tt.want {
				t.Errorf("got = %v, want = %v", got, tt.want)
			}
		})
	}
}
