package processjtl

import (
	"regexp"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

// MapToLuaTable converts a map[string]interface{} to a Lua table.
func MapToLuaTable(L *lua.LState, m map[string]interface{}) *lua.LTable {
	table := L.NewTable()

	for key, value := range m {
		switch key {
		case "Contents":
			// Make Contents accessible as .text
			table.RawSetString("text", lua.LString(value.(string)))
			// Also keep original Contents field
			table.RawSetString("Contents", lua.LString(value.(string)))
		case "class":
			// Make class accessible as .classname
			table.RawSetString("classname", lua.LString(value.(string)))
			// Also keep original class field
			table.RawSetString("class", lua.LString(value.(string)))
		case "id":
			// Make id accessible as .identifier
			table.RawSetString("identifier", lua.LString(value.(string)))
			// Also keep original id field
			table.RawSetString("id", lua.LString(value.(string)))
		case "children":
			if children, ok := value.([]interface{}); ok {
				childrenTable := L.NewTable()
				for i, child := range children {
					if childMap, ok := child.(map[string]interface{}); ok {
						childrenTable.RawSetInt(i+1, MapToLuaTable(L, childMap))
					}
				}
				table.RawSetString("children", childrenTable)
			}
		default:
			// Regular value conversion
			switch v := value.(type) {
			case string:
				table.RawSetString(key, lua.LString(v))
			case int:
				table.RawSetString(key, lua.LNumber(v))
			case float64:
				table.RawSetString(key, lua.LNumber(v))
			case bool:
				table.RawSetString(key, lua.LBool(v))
			case []interface{}:
				// Convert slice to Lua table
				arrayTable := L.NewTable()
				for i, item := range v {
					arrayTable.RawSetInt(i+1, convertToLuaValue(L, item))
				}
				table.RawSetString(key, arrayTable)
			case map[string]interface{}:
				// Recursively convert nested map
				table.RawSetString(key, MapToLuaTable(L, v))
			default:
				// Set nil for unsupported types
				table.RawSetString(key, lua.LNil)
			}
		}
	}

	return table
}

// convertToLuaValue converts an interface{} to a Lua value.
func convertToLuaValue(L *lua.LState, value interface{}) lua.LValue {
	switch v := value.(type) {
	case string:
		return lua.LString(v)
	case int:
		return lua.LNumber(v)
	case float64:
		return lua.LNumber(v)
	case bool:
		return lua.LBool(v)
	case []interface{}:
		arrayTable := L.NewTable()
		for i, item := range v {
			arrayTable.RawSetInt(i+1, convertToLuaValue(L, item))
		}
		return arrayTable
	case map[string]interface{}:
		return MapToLuaTable(L, v)
	default:
		return lua.LNil
	}
}

func luaTableToMap(table *lua.LTable) map[string]interface{} {
	result := make(map[string]interface{})
	table.ForEach(func(key, value lua.LValue) {
		switch v := value.(type) {
		case lua.LString:
			result[key.String()] = string(v)
		case lua.LNumber:
			result[key.String()] = float64(v)
		case lua.LBool:
			result[key.String()] = bool(v)
		case *lua.LTable:
			if key.String() == "children" {
				var children []interface{}
				v.ForEach(func(_, childValue lua.LValue) {
					if childTable, ok := childValue.(*lua.LTable); ok {
						children = append(children, luaTableToMap(childTable))
					}
				})
				result["children"] = children
			} else {
				result[key.String()] = luaTableToMap(v)
			}
		}
	})
	return result
}

func ParseCSS(css string) map[string]string {
	result := make(map[string]string)

	// Regex to match CSS properties.
	propertyRegex := regexp.MustCompile(`(?m)([\w-]+)\s*:\s*([^;]+)\s*;?`)

	// Find all properties within the CSS string.
	propMatches := propertyRegex.FindAllStringSubmatch(css, -1)
	// dont delete range here it is there for a reason
	for _, propMatch := range propMatches {
		key := strings.TrimSpace(propMatch[1])
		value := strings.TrimSpace(propMatch[2])
		result[key] = value
	}

	return result
}
