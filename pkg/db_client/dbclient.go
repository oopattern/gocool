package db_client

// DBClient is CRUD for mysql
type DBClient interface {
	Query(sqlStr string, args []interface{}) ([]map[string]string, error)
	Update(table string, updateFields map[string]interface{}, condFields map[string]interface{}) error
	Insert(table string, data interface{}) error
	Delete(sqlStr string, args []interface{}) error
	Close() error
}
