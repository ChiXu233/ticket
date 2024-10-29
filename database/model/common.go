package model

import (
	"database/sql/driver"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

var (
	offset   = 0
	limit    = 1
	OneQuery = QueryParams{
		Offset: &offset,
		Limit:  &limit,
	}
	EmptyFilter = make(map[string]interface{})
)

type ComparisonOperator string

const (
	NE   ComparisonOperator = "!="
	GT   ComparisonOperator = ">"
	LT   ComparisonOperator = "<"
	GE   ComparisonOperator = ">="
	LE   ComparisonOperator = "<="
	EQ   ComparisonOperator = "="
	LIKE ComparisonOperator = "like"
)

type Model struct {
	ID        int            `json:"id" gorm:"primaryKey;column:id"`
	CreatedAt time.Time      `json:"created_time,omitempty" gorm:"column:created_at"`
	UpdatedAt time.Time      `json:"updated_time,omitempty" gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at"`
}

//@TODO pg中有此问题，mysql并未发现。
// 问题描述：使用gorm时，结构体内使用的字段类型是“time.Time”
//// （1）查询返回结果为“2022-07-03T22:14:02.973528_08:00”，但我们需要“2022-07-03 22:14:02”这样的格式
//// （2）当时间字段不赋值时，插入到数据库会是“0001-01-01 00：00：00.000000+00：00”,我们不插入默认值
//// 解决方案：重新定义一个时间类型，并重写MarshalJson方法实现数据解析,重写Value和Scan方法实现存取数据时的相关操作

type LocalTime time.Time

const localTimeFormat = "2006-01-02 15:04:05"

func (t *LocalTime) UnmarshalJSON(data []byte) error {
	var err error
	if string(data) == "null" {
		return nil
	}

	now, err := time.ParseInLocation(`"`+localTimeFormat+`"`, string(data), time.Local)
	if err != nil {
		return err
	} else {
		*t = LocalTime(now)
	}
	return err
}

func (t LocalTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(localTimeFormat)+2)
	b = append(b, '"')
	if time.Time(t).Year() > 1000 {
		b = time.Time(t).AppendFormat(b, localTimeFormat)
		b = append(b, '"')
		return b, nil
	} else {
		b = append(b, '"')
		return b, nil
	}
}

func (t LocalTime) String() string {
	if time.Time(t).IsZero() {
		return "0000-00-00 00:00:00"
	}
	return time.Time(t).Format(localTimeFormat)
}

func (t LocalTime) Value() (driver.Value, error) {
	if time.Time(t).IsZero() {
		return nil, nil
	}
	return time.Time(t), nil
}

func (t *LocalTime) Scan(v interface{}) error {
	switch vt := v.(type) {
	case time.Time:
		*t = LocalTime(vt)
	case string:
		tTime, _ := time.Parse("2006/01/02 15:04:05", vt)
		*t = LocalTime(tTime)
	default:
		return nil
	}
	return nil
}

type QueryParams struct {
	SubQueries      []*SubQuery
	Orders          []Order // Orders
	Limit           *int    // Limit
	Offset          *int    // Offset
	InQueries       []*InQuery
	RangeQueries    []*RangeQuery
	CompareQueries  []*CompareQuery
	NotInQueries    []map[string]interface{}
	GroupFields     []string
	DistinctQueries []string
}

type CompareQuery struct {
	Field              string
	ComparisonOperator ComparisonOperator
	Value              interface{}
}

func (query *CompareQuery) Statement() string {
	if query.Field == "" || query.ComparisonOperator == "" || query.Value == nil {
		return ""
	}

	return fmt.Sprintf("%s %s ?", query.Field, query.ComparisonOperator)
}

func (query *CompareQuery) Parameters() []interface{} {
	parameters := make([]interface{}, 0, 1)
	parameters = append(parameters, query.Value)
	return parameters
}

type InQuery struct {
	Field  string
	Values interface{}
}

func (query *InQuery) Statement() string {
	if query.Field == "" || query.Values == nil {
		return ""
	}

	return query.Field + " IN (?)"
}

func (query *InQuery) Parameters() interface{} {
	return query.Values
}

type RangeQuery struct {
	Field string
	Start interface{}
	End   interface{}
}

func (query *RangeQuery) Statement() string {
	if query.Field == "" || query.Start == nil || query.End == nil {
		return ""
	}

	return query.Field + " BETWEEN ? AND ?"
}

func (query *RangeQuery) Parameters() []interface{} {
	parameters := make([]interface{}, 0, 2)
	parameters = append(parameters, query.Start)
	parameters = append(parameters, query.End)
	return parameters
}

type SubQuery struct {
	Keywords     []Keyword      // Keyword Filters, corresponding to `like` in sql select
	OrConditions []CompareQuery // or conditions, every map will be a group with 'AND'
}

func (query *SubQuery) Statement() string {
	if len(query.Keywords) == 0 {
		return ""
	}

	statements := make([]string, 0)
	for _, keyword := range query.Keywords {
		statements = append(statements, keyword.Field+" LIKE ?")
	}
	return strings.Join(statements, " OR ")
}

func (query *SubQuery) Parameters() []interface{} {
	parameters := make([]interface{}, 0)
	for _, keyword := range query.Keywords {
		switch keyword.Type {
		// 左模糊
		case 1:
			parameters = append(parameters, "%"+keyword.Value)
		// 右模糊
		case 2:
			parameters = append(parameters, keyword.Value+"%")
		default:
			parameters = append(parameters, "%"+keyword.Value+"%")
		}

	}
	return parameters
}

func (query *SubQuery) OrConstruct() string {
	if len(query.OrConditions) == 0 {
		return ""
	}

	statements := make([]string, 0)
	for i := range query.OrConditions {
		compareQuery := query.OrConditions[i]
		// bool、string做特殊处理
		value := convertParamToDBParam(compareQuery.Value)
		statements = append(statements, fmt.Sprintf("%s %s %v", compareQuery.Field, compareQuery.ComparisonOperator, value))

	}
	return strings.Join(statements, " OR ")
}

type Keyword struct {
	Field string
	Value string
	// 模糊查询是 左右/左/右/还是 0 1 2
	Type int
}

type Order struct {
	Field     string
	Direction string
}

type Pagination struct {
	Offset int `json:"offset"`
	Size   int `json:"size"`
}

func convertParamToDBParam(p interface{}) interface{} {
	switch p.(type) {
	case bool:
		if p.(bool) {
			return 1
		}
		return 0
	case string:
		if v, ok := p.(string); ok {
			return fmt.Sprintf("'%s'", v)
		}
		return p
	default:
		return p
	}
}
