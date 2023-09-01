package auth

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

/**
 *
 */
type PostgresAuth struct {
	pool        *pgxpool.Pool
	hash_option string
}

func initPostgresAuth(pgx *pgxpool.Pool, hash_option string) (*PostgresAuth, error) {
	return &PostgresAuth{
		pool:        pgx,
		hash_option: hash_option,
	}, nil

}

func (pa PostgresAuth) Login(Username string, Password string) string {

	conn, err := pa.pool.Acquire(context.Background())

	defer conn.Conn().Close(context.Background())

	if err != nil {
		log.Fatalf("Postgres Auth Can No Longer Acquire connections : %v", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	tx, errtx := conn.BeginTx(ctx, pgx.TxOptions{})

	sql := `
		SELECT password
		FROM Users
		WHERE username = @username
	`

	rows, err := tx.Query(ctx, sql, pgx.NamedArgs{
		"username": Username,
	})

	tx.Commit(ctx)

	defer rows.Close()

	if errtx == context.DeadlineExceeded {

	}

	if err == pgx.ErrNoRows {
		return ""
	}

	var password []byte

	/**
	 * or for rows.Next()
	 */

	_, err = pgx.ForEachRow(rows, []any{&password}, func() error {
		return nil
	})

	err = bcrypt.CompareHashAndPassword(password, []byte(password))

	if err != nil {
		return ""
	}

	return Username

}
