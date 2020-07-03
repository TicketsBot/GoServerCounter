package database

import(
	"database/sql"
	"fmt"
	"github.com/TicketsBot/GoServerCounter/config"
	_ "github.com/lib/pq"
)

type Database struct {
	*sql.DB
}

func NewDatabase() *Database {
	db, err := sql.Open("postgres", config.Conf.PostgresUri); if err != nil {
		panic(err)
	}

	return &Database{
		db,
	}
}

func (d *Database) GetServerCount() int {
	var count int

	if err := d.QueryRow(`SELECT COUNT(guild_id) FROM guilds;`).Scan(&count); err != nil {
		fmt.Printf("An error occurred whilst reading the server count: %s\n", err.Error())
		return 0
	}

	return count
}
