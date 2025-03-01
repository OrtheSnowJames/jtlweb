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

func ParseCSS(css string) map[string]string {
	result := make(map[string]string)

	// Regex to match CSS properties.
	propertyRegex := regexp.MustCompile(`(?m)([\w-]+)\s*:\s*([^;]+)\s*;?`)

	// Find all properties within the CSS string.
	propMatches := propertyRegex.FindAllStringSubmatch(css, -1)
	for _, propMatch := range propMatches {
		key := strings.TrimSpace(propMatch[1])
		value := strings.TrimSpace(propMatch[2])
		result[key] = value
	}

	return result
}
