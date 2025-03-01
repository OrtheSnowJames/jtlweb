package processjtl

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/OrtheSnowJames/jtl"
	lua "github.com/yuin/gopher-lua"
)

var document []interface{}
var documentUpdate atomic.Bool
var documentMutex sync.RWMutex
var ObjectsMutex sync.Mutex
var luaState *lua.LState

// Custom mutex-like Locker
type Locker struct {
	object interface{}
	lock   chan struct{}
}

func newLocker(obj interface{}) *Locker {
	return &Locker{
		object: obj,
		lock:   make(chan struct{}, 1),
	}
}

func (l *Locker) Lock() interface{} {
	l.lock <- struct{}{}
	return l.object
}

func (l *Locker) Unlock(obj interface{}) {
	l.object = obj
	<-l.lock
}

// getDocumentElement retrieves an element by class, ID, or key
func getDocumentElement(L *lua.LState) int {
	id := L.ToString(1)

	for _, elem := range document {
		elemMap, ok := elem.(map[string]interface{})
		if !ok {
			continue
		}

		if strings.HasPrefix(id, ".") {
			if class, ok := elemMap["class"].(string); ok && class == id[1:] {
				L.Push(MapToLuaTable(L, elemMap))
				return 1
			}
		}

		if strings.HasPrefix(id, "#") {
			if elemID, ok := elemMap["id"].(string); ok && elemID == id[1:] {
				L.Push(MapToLuaTable(L, elemMap))
				return 1
			}
		}

		if key, ok := elemMap["KEY"].(string); ok && key == id {
			L.Push(MapToLuaTable(L, elemMap))
			return 1
		}
	}

	L.Push(lua.LNil)
	return 1
}

// getObjects retrieves the objects as a Lua table
func getObjects(L *lua.LState) int {
	ObjectsMutex.Lock()
	defer ObjectsMutex.Unlock()

	objectsTable := L.NewTable()
	for i, obj := range objects {
		objTable := L.NewTable()
		objTable.RawSetString("type", lua.LString(reflect.TypeOf(obj).String()))
		objTable.RawSetString("index", lua.LNumber(i))
		objectsTable.Append(objTable)
	}

	L.Push(objectsTable)
	return 1
}

// Extract all script contents from JTL document
func extractScripts(jtlcomps []interface{}) string {
	var scripts strings.Builder
	for _, elem := range jtlcomps {
		comp, ok := elem.(map[string]interface{})
		if !ok {
			continue
		}

		if key, ok := comp["KEY"].(string); ok && key == "script" {
			if content, ok := comp["Contents"].(string); ok {
				scripts.WriteString(content)
				scripts.WriteString("\n")
			}
		}
	}
	return scripts.String()
}

// updateDocument safely updates the document
func updateDocument(L *lua.LState) int {
	newDoc := L.ToTable(1)
	if newDoc == nil {
		return 0
	}

	var docArray []interface{}
	newDoc.ForEach(func(idx, value lua.LValue) {
		if tbl, ok := value.(*lua.LTable); ok {
			docArray = append(docArray, luaTableToMap(tbl))
		}
	})

	documentMutex.Lock()
	document = docArray
	documentMutex.Unlock()
	documentUpdate.Store(true)

	return 0
}

func getElement(selector string) UIElement {
	for _, obj := range objects {
		if baseEl, ok := obj.(interface{ GetBaseElement() *BaseElement }); ok {
			el := baseEl.GetBaseElement()
			if strings.HasPrefix(selector, ".") && el.Class == selector[1:] {
				return obj.(UIElement)
			}
			if strings.HasPrefix(selector, "#") && el.ID == selector[1:] {
				return obj.(UIElement)
			}
		}
	}
	return nil
}

// Add new Lua functions
func setEventHandler(L *lua.LState) int {
	selector := L.ToString(1)
	event := L.ToString(2)
	handler := L.ToString(3)

	if element := getElement(selector); element != nil {
		if baseEl, ok := element.(interface{ GetBaseElement() *BaseElement }); ok {
			baseEl.GetBaseElement().SetEventHandler(event, handler)
		}
	}
	return 0
}

func executeEventHandler(element *BaseElement, event string) {
	if handler := element.GetEventHandler(event); handler != "" && luaState != nil {
		if err := luaState.DoString(handler); err != nil {
			fmt.Printf("Error executing event handler: %v\n", err)
		}
	}
}

// MakeWebview now prepares view without creating a new window
func MakeWebview(jtldoc string) (*Locker, []CanvasObject) {
	luaState = lua.NewState()
	defer luaState.Close()

	// Parse JTL document
	parsedDoc, err := jtl.Parse(jtldoc)
	if err != nil {
		fmt.Printf("Failed to parse JTL: %v\n", err)
		return nil, nil
	}

	// Extract scripts
	combinedScript := extractScripts(parsedDoc)

	// Setup Lua environment
	docTable := luaState.NewTable()
	luaState.SetField(docTable, "get", luaState.NewFunction(getDocumentElement))
	luaState.SetField(docTable, "objects", luaState.NewFunction(getObjects))
	luaState.SetField(docTable, "update", luaState.NewFunction(updateDocument))
	luaState.SetGlobal("document", docTable)

	// Add event handling function
	luaState.SetGlobal("onEvent", luaState.NewFunction(setEventHandler))

	documentMutex.Lock()
	document = parsedDoc
	documentMutex.Unlock()
	objects := ToRaylib(document)

	// Background goroutine for script execution and document updates
	go func() {
		scriptL := lua.NewState()
		defer scriptL.Close()

		// Setup same environment for script execution
		scriptDocTable := scriptL.NewTable()
		scriptL.SetField(scriptDocTable, "get", scriptL.NewFunction(getDocumentElement))
		scriptL.SetField(scriptDocTable, "objects", scriptL.NewFunction(getObjects))
		scriptL.SetField(scriptDocTable, "update", scriptL.NewFunction(updateDocument))
		scriptL.SetGlobal("document", scriptDocTable)

		// Execute initial script
		if err := scriptL.DoString(combinedScript); err != nil {
			fmt.Printf("Script execution error: %v\n", err)
		}

		// Watch for document updates
		for {
			time.Sleep(1 * time.Second)
			if documentUpdate.Load() {
				documentUpdate.Store(false)
				documentMutex.RLock()
				newObjects := ToRaylib(document)
				documentMutex.RUnlock()

				ObjectsMutex.Lock()
				objects = newObjects
				ObjectsMutex.Unlock()

				// Re-run script after document update
				if err := scriptL.DoString(combinedScript); err != nil {
					fmt.Printf("Script re-execution error: %v\n", err)
				}
			}
		}
	}()

	return newLocker(objects), objects
}

// Helper function to convert Lua table to Go map
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
			result[key.String()] = luaTableToMap(v)
		}
	})
	return result
}
