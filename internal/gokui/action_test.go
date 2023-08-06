package gokui

import "testing"

func TestGenerateSelect(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want string
	}{
		{
			sql: `
CREATE TABLE users (
    user_id VARCHAR(256) NOT NULL COMMENT 'user id',
    name VARCHAR(256),
    age INT UNSIGNED NOT NULL,
    created_at DATETIME(6) NOT NULL,
    primary key (user_id)
);
`,
			want: `SELECT user_id, name, age, created_at FROM users;`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateSelect(tt.sql)
			if err != nil {
				t.Errorf("error occurred: %v", err)
			}
			if got != tt.want {
				t.Errorf("got = %v, want = %v", got, tt.want)
			}
		})
	}
}
