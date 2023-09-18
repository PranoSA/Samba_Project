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

func InitPostgresAuth(pgx *pgxpool.Pool, hash_option string) (*PostgresAuth, error) {
	return &PostgresAuth{
		pool:        pgx,
		hash_option: hash_option,
	}, nil

}

func (pa PostgresAuth) Login(Username string, Password string) bool {

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
		return false
	}

	var password []byte

	/**
	 * or for rows.Next()
	 */

	_, err = pgx.ForEachRow(rows, []any{&password}, func() error {
		return nil
	})

	if pa.hash_option == "bcrypt" {

		err = bcrypt.CompareHashAndPassword(password, []byte(password))
	}

	if pa.hash_option == "md5" {
		log.Fatalf("MD5 Hash Option Not Supported Yet")
	}

	if pa.hash_option == "argon" {
		log.Fatal("Argon Hash Option Not Supported Yet")
	}

	if err != nil {
		return false
	}

	return true

}

func (pa PostgresAuth) Signup(Username string, Password string) bool {
	return true
}
