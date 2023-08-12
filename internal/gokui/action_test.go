package gokui

import "testing"

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
		name string
		sql  string
		opt  GenerateSelectOptions
		want string
	}{
		{
			name: "newline = false",
			sql:  createTableUsers,
			opt: GenerateSelectOptions{
				NewLine: false,
			},
			want: `SELECT user_id, name, age, created_at FROM users;`,
		},
		{
			name: "newline = true",
			sql:  createTableUsers,
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
		t.Run(tt.name, func(t *testing.T) {
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
		name string
		sql  string
		opt  GenerateInsertOptions
		want string
	}{
		{
			name: "newline = false, insert-select = false",
			sql:  createTableUsers,
			opt: GenerateInsertOptions{
				NewLine:      false,
				InsertSelect: false,
			},
			want: `INSERT INTO users (user_id, name, age, created_at) VALUES ('', '', 0, '');`,
		},
		{
			name: "newline = true, insert-select = false",
			sql:  createTableUsers,
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
			name: "newline = false, insert-select = true",
			sql:  createTableUsers,
			opt: GenerateInsertOptions{
				NewLine:      false,
				InsertSelect: true,
			},
			want: `INSERT INTO users (user_id, name, age, created_at) SELECT user_id, name, age, created_at FROM users;`,
		},
		{
			name: "newline = true, insert-select = true",
			sql:  createTableUsers,
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
		t.Run(tt.name, func(t *testing.T) {
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
