package sqlgo

import (
	"fmt"
	"strconv"
)

type Serializer struct {
	params []interface{}
}

func (s *Serializer) Add(p interface{}) string {
	i := len(s.params) + 1
	s.params = append(s.params, p)
	return fmt.Sprintf("$%v", i)
}

func (s *Serializer) Params() []interface{} {
	return s.params
}

func NewSerializer() *Serializer {
	s := Serializer{}
	s.params = make([]interface{}, 0)
	return &s
}

// Serializable defines the interface of types that can be written
// into the database.
type Serializable interface {
	// GenerateInsertSQL generates the SQL to use when inserting
	// the serializable into the database
	GenerateInsertSQL() string
}

// Serialize serializes the given interface{} into a SQL string
func Serialize(s interface{}) string {
	if s == nil {
		return "NULL"
	}

	switch v := s.(type) {
	case string:
		return "'" + v + "'"
	case int:
		return strconv.FormatInt(int64(v), 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		if v {
			return "true"
		}
		return "false"
	default:
		panic("Tried to serialize unknown type!")
	}
}

// SerializeStringArray serializes an array of strings into SQL
func SerializeStringArray(s []string) string {
	result := "ARRAY["

	for i, item := range s {
		result += Serialize(item)

		if i != len(s)-1 {
			result += ","
		}
	}

	result += "]"

	return result
}
