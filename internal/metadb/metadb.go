package metadb

// MetaDB is an interface to relational databases.
type MetaDB interface {
	Query(query string, arg interface{}, result interface{}) error
	Execute(query string, arg interface{}) error
}
