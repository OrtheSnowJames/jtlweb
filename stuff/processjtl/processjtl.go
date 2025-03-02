package processjtl

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"

	"jtlweb/stuff/shared"

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

var documentStore []map[string]interface{}

func clearDocuments() {
	documentStore = make([]map[string]interface{}, 0)
	updateObjectsFromDocumentStore()
}

func insertDocument(attributes map[string]interface{}) {
	documentStore = append(documentStore, attributes)
	updateObjectsFromDocumentStore()
}

func getDocumentsByAttribute(key string, value interface{}) []map[string]interface{} {
	var results []map[string]interface{}
	for _, doc := range documentStore {
		if docValue, exists := doc[key]; exists && fmt.Sprint(docValue) == fmt.Sprint(value) {
			results = append(results, doc)
		}
	}
	return results
}

func removeDocumentByAttribute(key string, value interface{}) {
	newDocs := make([]map[string]interface{}, 0)
	for _, doc := range documentStore {
		if docValue, exists := doc[key]; !exists || fmt.Sprint(docValue) != fmt.Sprint(value) {
			newDocs = append(newDocs, doc)
		}
	}
	documentStore = newDocs
	updateObjectsFromDocumentStore()
}

func getAllDocuments() []interface{} {
	result := make([]interface{}, len(documentStore))
	for i, doc := range documentStore {
		result[i] = doc
	}
	return result
}

func updateObjectsFromDocumentStore() {
	documentMutex.Lock()
	document = getAllDocuments()
	documentMutex.Unlock()
	newObjects := ToRaylib(document)
	ObjectsMutex.Lock()
	objects = newObjects
	ObjectsMutex.Unlock()
	fmt.Printf("Objects updated, new length: %d\n", len(objects))
}

// getDocumentElement retrieves an element by class, ID, or key
func getDocumentElement(L *lua.LState) int {
	id := L.ToString(1)
	fmt.Printf("Searching for element with id: %s\n", id)

	var docs []map[string]interface{}
	var searchKey string
	var searchVal interface{}

	if strings.HasPrefix(id, ".") {
		searchKey = "class"
		searchVal = id[1:]
	} else if strings.HasPrefix(id, "#") {
		searchKey = "id"
		searchVal = id[1:]
	} else {
		searchKey = "KEY"
		searchVal = id
	}

	docs = getDocumentsByAttribute(searchKey, searchVal)
	if len(docs) == 0 {
		L.Push(lua.LNil)
		return 1
	}

	// Convert all document attributes directly to Lua table
	table := MapToLuaTable(L, docs[0])

	// Immediately update objects on removal:
	removeFunc := L.NewFunction(func(L *lua.LState) int {
		removeDocumentByAttribute(searchKey, searchVal)
		updateObjectsFromDocumentStore() // Ensure objects are updated immediately
		return 0
	})
	table.RawSetString("remove", removeFunc)

	L.Push(table)
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

// Will be updated later for jtltp (connection) support
func GetRelToOpenPath(relpath string) string {
	// Get the directory containing the current JTL file
	dir := filepath.Dir(shared.OpenPath)
	// Join the directory with the relative path
	return filepath.Join(dir, relpath)
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
			if content, ok := comp["Contents"].(string); ok && (content != "" && content != "\n") {
				scripts.WriteString(content)
				scripts.WriteString("\n")
			} else if relpath, ok := comp["src"].(string); ok {
				path := GetRelToOpenPath(relpath)
				fmt.Printf("path: %v\n", path)
				script, err := os.ReadFile(path)
				if err != nil {
					fmt.Printf("Error reading script file: %v\n", err)
					continue
				}

				scripts.WriteString(string(script))
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
	updateObjectsFromDocumentStore() // Ensure objects are updated immediately

	return 0
}

func getElement(selector string) UIElement {
	ObjectsMutex.Lock()
	defer ObjectsMutex.Unlock()

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

// Helper function to setup Lua environment
func setupLuaEnvironment(L *lua.LState) *lua.LTable {
	docTable := L.NewTable()
	L.SetField(docTable, "get", L.NewFunction(getDocumentElement))
	L.SetField(docTable, "objects", L.NewFunction(getObjects))
	L.SetField(docTable, "update", L.NewFunction(updateDocument))
	L.SetField(docTable, "onEvent", L.NewFunction(setEventHandler))
	return docTable
}

// MakeWebview now prepares view without creating a new window
func MakeWebview(jtldoc string) (*Locker, []CanvasObject) {
	luaState = lua.NewState()

	// Parse JTL document
	parsedDoc, err := jtl.Parse(jtldoc)
	if err != nil {
		fmt.Printf("Failed to parse JTL: %v\n", err)
		return nil, nil
	}

	fmt.Printf("Parsed %d JTL components\n", len(parsedDoc))

	// Clear existing documents
	clearDocuments()

	// Store in memory
	for _, elem := range parsedDoc {
		if elemMap, ok := elem.(map[string]interface{}); ok {
			insertDocument(elemMap)
		}
	}

	// Create objects from all documents
	allDocs := getAllDocuments()
	if len(allDocs) == 0 {
		fmt.Println("No documents retrieved from database")
		return nil, nil
	}

	objects = ToRaylib(allDocs)
	if len(objects) == 0 {
		fmt.Println("No objects created from documents")
		return nil, nil
	}

	fmt.Printf("Created %d objects\n", len(objects))

	// Extract and run scripts after objects are created
	combinedScript := extractScripts(parsedDoc)

	// Setup initial Lua environment
	docTable := setupLuaEnvironment(luaState)
	luaState.SetGlobal("document", docTable)

	// Execute script after objects are created and stored
	if err := luaState.DoString(combinedScript); err != nil {
		fmt.Printf("Initial script execution error: %v\n", err)
	}

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
