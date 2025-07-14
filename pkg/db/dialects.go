package db

type Dialect string

func (d Dialect) String() string {
	return string(d)
}

type ConnInfo struct {
	Dialect  Dialect
	Host     string
	Port     string
	User     string
	Database string
	Password string
}

const (
	Postgres Dialect = "postgres"
	MySQL    Dialect = "mysql"
)
