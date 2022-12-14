package schema

import (
	"gee/dialect"
	"go/ast"
	"reflect"
)

type Field struct {
	Name string
	Type string
	Tag  string
}

// Schema represents a table of database
type Schema struct {
	Model      interface{} //model对象
	Name       string
	Fields     []*Field
	FieldNames []string
	fieldMap   map[string]*Field
}

func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

func Parse(dest interface{}, d dialect.Dialect) *Schema {
	//TypeOf()和ValueOf分别是返回入参的类型和值
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(), //获取结构体的名称作为表名
		fieldMap: make(map[string]*Field),
	}
	//NumField()获取实例的字段的个数
	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i) //通过下标获取特定的字段
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,                                              //字段名
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))), //将p.type转换为数据库的字段类型
			}
			//p.Tag额外的约束条件
			if v, ok := p.Tag.Lookup("geeorm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}
