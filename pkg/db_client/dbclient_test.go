package db_client

import (
	"testing"
)

func TestBuildGrpcServer(t *testing.T) {
	defaultCfg := DBConfig{
		DBType: "test",
		DBHost: "127.0.0.1",
		DBPort: 3306,
		DBPass: "xxx",
		DBUser: "root",
		DBName: "test_db",
	}

	dbClient, err := NewMysqlClient(defaultCfg)
	if err != nil {
		t.Error(err)
	}

	defer dbClient.Close()

	sql := "select count(1) from test_table limit 1"
	results, err := dbClient.Query(sql, nil)
	if err != nil {
		t.Error(err)
	}

	t.Logf("size=%d", len(results))
}
