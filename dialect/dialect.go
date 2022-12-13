package dialect

import "reflect"

var dialectsMap = map[string]Dialect{}

type Dialect interface {
	DataTypeOf(typ reflect.Value) string	//转化为数据库的数据类型
	TableExistSQL(tableName string) (string, []interface{})	//表是否存在
}

func RegisterDialect(name string, dialect Dialect){
	dialectsMap[name] = dialect
}

func GetDialect(name string) (dialect Dialect, ok bool){
	dialect, ok = dialectsMap[name]
	return
}