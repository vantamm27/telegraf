package json

import (
	"fmt"
	"strconv"
)

type JSONFlattener struct {
	Fields map[string]interface{}
}

// FlattenJSON flattens nested maps/interfaces into a fields map (ignoring bools and string)
func (f *JSONFlattener) FlattenJSON(
	fieldname string,
	v interface{}) error {
	//("json", "FlattenJSON")
	if f.Fields == nil {
		f.Fields = make(map[string]interface{})
	}

	return f.FullFlattenJSON(fieldname, v, false, false, 0)
}

// FullFlattenJSON flattens nested maps/interfaces into a fields map (including bools and string)
func (f *JSONFlattener) FullFlattenJSON(
	fieldname string,
	v interface{},
	convertString bool,
	convertBool bool,
	dep int,
) error {
	//fmt.Println("json", "FullFlattenJSON", dep)
	//fmt.Println(fieldname, fieldname, convertString, convertBool)
	if f.Fields == nil {
		f.Fields = make(map[string]interface{})
	}

	switch t := v.(type) {
	case map[string]interface{}:
		if dep >= 1 {
			f.Fields[fieldname] = t
		} else {

			for k, v := range t {
				fieldkey := k
				if fieldname != "" {
					fieldkey = fieldname + "_" + fieldkey
				}

				err := f.FullFlattenJSON(fieldkey, v, convertString, convertBool, dep+1)
				if err != nil {
					return err
				}
			}
		}
		//f.Fields[fieldname] = t
	case []interface{}:
		if dep >= 1 {
			f.Fields[fieldname] = t
		} else {
			for i, v := range t {
				fieldkey := strconv.Itoa(i)
				if fieldname != "" {
					fieldkey = fieldname + "_" + fieldkey
				}
				err := f.FullFlattenJSON(fieldkey, v, convertString, convertBool, dep+1)
				if err != nil {
					return err
				}
			}
		}
		// f.Fields[fieldname] = t
	case float64:
		f.Fields[fieldname] = t
	case string:
		if !convertString {
			return nil
		}
		f.Fields[fieldname] = v.(string)
	case bool:
		if !convertBool {
			return nil
		}
		f.Fields[fieldname] = v.(bool)
	case nil:
		return nil
	default:
		return fmt.Errorf("JSON Flattener: got unexpected type %T with value %v (%s)",
			t, t, fieldname)
	}
	return nil
}
