package api

import (
	"authencation-service/src/data"
	"database/sql"
)

type Config struct {
	WebPort int
	DB      *sql.DB
	Models  data.Models
}
