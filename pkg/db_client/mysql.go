package db_client

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/oopattern/gocool/log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var defaultDBOptions = &DBOptions{
	MaxLiftTime: 15,
	MaxIdleConn: 1000,
	MaxOpenConn: 1000,
}

type MysqlClient struct {
	DBHandler interface{}
	Tx        bool
}

func NewMysqlClient(cfg DBConfig) (DBClient, error) {
	conn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)

	handler, err := sql.Open("mysql", conn)
	if err != nil {
		log.Error("init default client error: %v", err)
		return nil, err
	}

	// 15分钟链接没有使用就断开
	handler.SetConnMaxLifetime(time.Duration(defaultDBOptions.MaxLiftTime) * time.Minute)
	handler.SetMaxIdleConns(defaultDBOptions.MaxIdleConn)
	handler.SetMaxOpenConns(defaultDBOptions.MaxOpenConn)

	return &MysqlClient{
		DBHandler: handler,
		Tx:        false}, nil
}

func (c *MysqlClient) Close() error {
	return c.DBHandler.(*sql.DB).Close()
}

func (c *MysqlClient) Query(sqlStr string, args []interface{}) ([]map[string]string, error) {
	debugSql := strings.Replace(sqlStr, "?", "%v", -1)
	debugSql = fmt.Sprintf(debugSql, args...)
	log.Debug("query sql = %s", debugSql)

	rows, err := c.DBHandler.(*sql.DB).Query(sqlStr, args...)
	if err != nil {
		log.Error("query db error, err = %s, sql = %s, args = %v, debug_sql = %s", err.Error(), sqlStr, args, debugSql)
		return nil, err
	}

	result, err := c.parseRows(rows)
	if err != nil {
		log.Error("parse rows error, err = %v", err)
		return nil, err
	}

	return result, nil
}

func (c *MysqlClient) Update(table string, updateFields map[string]interface{}, condFields map[string]interface{}) error {
	if 0 == len(table) || 0 == len(updateFields) || 0 == len(condFields) {
		return errors.New("field map empty")
	}

	sqlStr := fmt.Sprintf("update %s set ", table)
	args := make([]interface{}, 0)
	for k, v := range updateFields {
		sqlStr += fmt.Sprintf("%s = ?, ", k)
		args = append(args, v)
	}
	sqlStr = strings.TrimRight(sqlStr, ", ")

	sqlStr += " where "
	for k, v := range condFields {
		sqlStr += fmt.Sprintf("%s = ? and ", k)
		args = append(args, v)
	}
	sqlStr = strings.TrimRight(sqlStr, "and ")

	size, err := c.execDB(sqlStr, args)
	if err != nil {
		log.Error("ExecDb error, err = %v", err)
		return err
	}

	log.Info("update data count: %v", size)
	return nil
}

func (c *MysqlClient) Insert(table string, data interface{}) error {
	// 拼接sql
	sqlStr := fmt.Sprintf("insert into %s (", table)
	args := make([]interface{}, 0)
	typeFields := reflect.TypeOf(data)
	valueFields := reflect.ValueOf(data)
	for i := 0; i < typeFields.NumField(); i++ {
		field := typeFields.Field(i)
		sqlStr += field.Tag.Get("json") + ", "
		value := valueFields.Field(i).Interface()
		args = append(args, value)
	}
	sqlStr = strings.TrimRight(sqlStr, ", ")
	sqlStr += ") values (?" + strings.Repeat(",?", typeFields.NumField()-1) + ")"

	// 执行sql
	row_size, err := c.execDB(sqlStr, args)
	if err != nil {
		log.Error("ExecDb error, err = %v", err)
		return err
	}

	log.Info("affected row size = %d", row_size)
	return nil
}

func (c *MysqlClient) Delete(sqlStr string, args []interface{}) error {
	return nil
}

func (c *MysqlClient) execDB(sqlStr string, args []interface{}) (int64, error) {
	debugSql := strings.Replace(sqlStr, "?", "%v", -1)
	debugSql = fmt.Sprintf(debugSql, args...)
	log.Debug("exec sql = %s", debugSql)

	var rowSize int64
	if c.Tx {
		txDbHandler := c.DBHandler.(*sql.Tx)
		result, err := txDbHandler.Exec(sqlStr, args...)
		if err != nil {
			log.Error("exec error, err = %s, sql = %s, args = %v", err.Error(), sqlStr, args)
			return 0, err
		}
		rowSize, err = result.RowsAffected()
		if err != nil {
			log.Error("decode affected row err: %v", err.Error(), sqlStr, args)
			return 0, err
		}
		log.Info("tx affected row size = %d", rowSize)
	} else {
		dbHandler := c.DBHandler.(*sql.DB)
		result, err := dbHandler.Exec(sqlStr, args...)
		if err != nil {
			log.Error("exec error, err = %s, sql = %s, args = %v", err.Error(), sqlStr, args)
			return 0, err
		}
		rowSize, err = result.RowsAffected()
		if err != nil {
			log.Error("decode affected row err: %v", err.Error(), sqlStr, args)
			return 0, err
		}
		log.Info("affected row size = %d", rowSize)
	}

	return rowSize, nil
}

func (c *MysqlClient) parseRows(rows *sql.Rows) ([]map[string]string, error) {
	defer func() {
		if err := rows.Close(); err != nil {
			log.Error("rows close err[%v]", err)
		}
	}()

	columns, _ := rows.Columns()
	columnTypes, _ := rows.ColumnTypes()
	arr := make([]interface{}, len(columns))
	for i, v := range columnTypes {
		t := v.ScanType()
		v := reflect.New(t).Interface()
		arr[i] = v
	}

	fullValues := make([]map[string]string, 0)
	for rows.Next() {
		err := rows.Scan(arr...)
		if err != nil {
			log.Error("row scan error, err = %v", err)
			return nil, err
		}
		records := make(map[string]string)
		for i, col := range columns {
			v := arr[i]
			switch vv := v.(type) {
			case *int32:
				records[col] = strconv.FormatInt(int64(*vv), 10)
			case *uint32:
				records[col] = strconv.FormatUint(uint64(*vv), 10)
			case *int64:
				records[col] = strconv.FormatInt(int64(*vv), 10)
			case *uint64:
				records[col] = strconv.FormatUint(uint64(*vv), 10)
			case *sql.NullString:
				if vv.Valid {
					records[col] = vv.String
				} else {
					records[col] = ""
				}
			case *sql.NullBool:
				if vv.Valid {
					if vv.Bool {
						records[col] = "true"
					} else {
						records[col] = "false"
					}
				} else {
					records[col] = "false"
				}
			case *sql.NullFloat64:
				if vv.Valid {
					records[col] = strconv.FormatFloat(vv.Float64, 'E', -1, 64)
				} else {
					records[col] = "0.0"
				}
			case *sql.NullInt64:
				if vv.Valid {
					records[col] = strconv.FormatInt(vv.Int64, 10)
				} else {
					records[col] = "0"
				}
			case *mysql.NullTime:
				if vv.Valid {
					records[col] = vv.Time.Format("2006-01-02 15:04:05")
				} else {
					records[col] = "2006-01-02 15:04:05"
				}
			case *sql.RawBytes:
				records[col] = string(*vv)
			default:
				records[col] = string(v.([]byte))
			}
		}
		fullValues = append(fullValues, records)
	}
	return fullValues, nil
}
