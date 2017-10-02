package data

type DB interface {
	Connect() error
	FindByHash(string) (*Avatar, error)
	Save(*Avatar) error
	Migrate() error
}
