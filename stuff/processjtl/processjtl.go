package processjtl

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"jtlweb/stuff/jtltp"
	"jtlweb/stuff/shared"

	"github.com/OrtheSnowJames/jtl"
	lua "github.com/yuin/gopher-lua"
)

var document []interface{}
var documentUpdate atomic.Bool
var documentMutex sync.RWMutex
var ObjectsMutex sync.Mutex
var luaState *lua.LState
var Site string
var frameHandler string          // Add this at the top with other vars
var requestedFrameHandler string // Add this at the top with other vars

var updateMainObjectsCallback func([]CanvasObject)

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
var createdElements []map[string]interface{} // Store created elements

func clearDocuments() {
	documentStore = make([]map[string]interface{}, 0)
	updateObjectsFromDocumentStore()
}

func insertDocument(attributes map[string]interface{}) {
	documentStore = append(documentStore, attributes)
	updateObjectsFromDocumentStore()
}

func insertDocumentAt(index int, attributes map[string]interface{}) {
	if index < 0 || index > len(documentStore) {
		index = len(documentStore)
	}
	documentStore = append(documentStore[:index], append([]map[string]interface{}{attributes}, documentStore[index:]...)...)
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

func replaceDocumentByAttribute(key string, value interface{}, newElement map[string]interface{}) {
	removeDocumentByAttribute(key, value)
	insertDocument(newElement)
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

	// Recalculate positions of all objects and update content height
	contentHeight := 0
	for _, obj := range objects {
		if baseEl, ok := obj.(interface{ GetBaseElement() *BaseElement }); ok {
			contentHeight += int(baseEl.GetBaseElement().Height) + 20 // Add height and margin
		}
	}
	shared.ContentHeight = contentHeight

	// Call the callback to update main objects
	if updateMainObjectsCallback != nil {
		updateMainObjectsCallback(objects)
	}
}

// SetUpdateMainObjectsCallback sets the callback function to update main objects
func SetUpdateMainObjectsCallback(callback func([]CanvasObject)) {
	updateMainObjectsCallback = callback
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

	// Retrieve the text property of a TextField
	if textField, ok := docs[0]["text"].(string); ok {
		table.RawSetString("text", lua.LString(textField))
	}

	L.Push(table)
	return 1
}

// GetDocumentElements retrieves all elements in a class or with the same key
func GetDocumentElements(classkey string) []map[string]interface{} {
	if strings.HasPrefix(classkey, ".") {
		return getDocumentsByAttribute("class", classkey[1:])
	} else {
		return getDocumentsByAttribute("KEY", classkey)
	}
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

// Add document.create, document.add, document.replace, document.insert, and document.whereis functions
func createElement(L *lua.LState) int {
	elementType := L.ToString(1)
	content := L.ToString(2)
	styles := L.ToTable(3)

	element := map[string]interface{}{
		"KEY":      elementType,
		"Contents": content,
	}

	if styles != nil {
		styles.ForEach(func(key, value lua.LValue) {
			element[key.String()] = value.String()
		})
	}

	createdElements = append(createdElements, element)
	L.Push(MapToLuaTable(L, element))
	return 1
}

// Fix x.text (in lua api) not changing to elem.Contents
func addElement(L *lua.LState) int {
	element := L.ToTable(1)
	if element == nil {
		return 0
	}

	// Convert Lua table to Go map
	newElement := luaTableToMap(element)

	// Ensure the Contents field is set correctly
	if text, ok := newElement["text"].(string); ok {
		newElement["Contents"] = text
	}

	// Ensure the KEY field is set correctly
	if key, ok := newElement["KEY"].(string); !ok || key == "" {
		newElement["KEY"] = "p" // Default to "p" if KEY is not set
		fmt.Println("KEY not set, defaulting to 'p'")
	}

	// json the element
	jsonElement, err := json.Marshal(newElement)
	if err != nil {
		fmt.Printf("Error marshalling element: %v\n", err)
		return 0
	}
	// string the json
	fmt.Printf("Adding element: %v\n", string(jsonElement))
	insertDocument(newElement)
	return 0
}

func insertElement(L *lua.LState) int {
	index := L.ToInt(1)
	element := L.ToTable(2)
	if element == nil {
		return 0
	}

	// Convert Lua table to Go map
	newElement := luaTableToMap(element)

	// Ensure the Contents field is set correctly
	if text, ok := newElement["text"].(string); ok {
		newElement["Contents"] = text
	}

	// Ensure the KEY field is set correctly
	if key, ok := newElement["KEY"].(string); !ok || key == "" {
		newElement["KEY"] = "p" // Default to "p" if KEY is not set
		fmt.Println("KEY not set, defaulting to 'p'")
	}

	// json the element
	jsonElement, err := json.Marshal(newElement)
	if err != nil {
		fmt.Printf("Error marshalling element: %v\n", err)
		return 0
	}
	// string the json
	fmt.Printf("Inserting element at index %d: %v\n", index, string(jsonElement))
	insertDocumentAt(index-1, newElement) // Lua index starts at 1
	return 0
}

func replaceElement(L *lua.LState) int {
	selector := L.ToString(1)
	newElement := L.ToTable(2)
	if newElement == nil {
		return 0
	}

	// Convert Lua table to Go map
	newElementMap := luaTableToMap(newElement)

	// Ensure the Contents field is set correctly
	if text, ok := newElementMap["text"].(string); ok {
		newElementMap["Contents"] = text
	}

	// Ensure the KEY field is set correctly
	if key, ok := newElementMap["KEY"].(string); !ok || key == "" {
		newElementMap["KEY"] = "p" // Default to "p" if KEY is not set
		fmt.Println("KEY not set, defaulting to 'p'")
	}

	var searchKey string
	var searchVal interface{}

	if strings.HasPrefix(selector, ".") {
		searchKey = "class"
		searchVal = selector[1:]
	} else if strings.HasPrefix(selector, "#") {
		searchKey = "id"
		searchVal = selector[1:]
	} else {
		searchKey = "KEY"
		searchVal = selector
	}

	// Find the index of the element to replace
	index := -1
	for i, doc := range documentStore {
		if docValue, exists := doc[searchKey]; exists && fmt.Sprint(docValue) == fmt.Sprint(searchVal) {
			index = i
			break
		}
	}

	// If the element is found, remove it and insert the new element at the same index
	if index != -1 {
		removeDocumentByAttribute(searchKey, searchVal)
		insertDocumentAt(index, newElementMap)
	} else {
		// If the element is not found, just insert the new element
		insertDocument(newElementMap)
	}

	updateObjectsFromDocumentStore()
	return 0
}

func whereisElement(L *lua.LState) int {
	selector := L.ToString(1)

	var searchKey string
	var searchVal interface{}

	if strings.HasPrefix(selector, ".") {
		searchKey = "class"
		searchVal = selector[1:]
	} else if strings.HasPrefix(selector, "#") {
		searchKey = "id"
		searchVal = selector[1:]
	} else {
		searchKey = "KEY"
		searchVal = selector
	}

	for i, doc := range documentStore {
		if docValue, exists := doc[searchKey]; exists && fmt.Sprint(docValue) == fmt.Sprint(searchVal) {
			L.Push(lua.LNumber(i + 1)) // Lua index starts at 1
			return 1
		}
	}

	L.Push(lua.LNil)
	return 1
}

// Helper function to setup Lua environment
func setupLuaEnvironment(L *lua.LState) *lua.LTable {
	docTable := L.NewTable()
	L.SetField(docTable, "get", L.NewFunction(getDocumentElement))
	L.SetField(docTable, "objects", L.NewFunction(getObjects))
	L.SetField(docTable, "update", L.NewFunction(updateDocument))
	L.SetField(docTable, "onEvent", L.NewFunction(setEventHandler))
	L.SetField(docTable, "create", L.NewFunction(createElement))
	L.SetField(docTable, "add", L.NewFunction(addElement))
	L.SetField(docTable, "replace", L.NewFunction(replaceElement))
	L.SetField(docTable, "insert", L.NewFunction(insertElement))
	L.SetField(docTable, "whereis", L.NewFunction(whereisElement))
	L.SetField(docTable, "addStyle", L.NewFunction(addStyle))
	L.SetField(docTable, "removeAllStyle", L.NewFunction(removeAllStyle))
	L.SetField(docTable, "fetch", L.NewFunction(luaWrapFetch))
	L.SetField(docTable, "onFrame", L.NewFunction(setFrameHandler))
	L.SetField(docTable, "requestFrame", L.NewFunction(setRequestFrameHandler))
	return docTable
}

func luaWrapFetch(L *lua.LState) int {
	// get what to get at ("here")
	selector := L.ToString(1)

	resp, err := jtltp.JtltpFetch(Site, selector)
	if err != nil {
		fmt.Printf("Error fetching: %v\n", err)
		return 0
	}

	// convert to lua table: has {"JTLTP-STATUS": 200 (hopefully), "JTLTP-TYPE": "jtl", "JTLTP": ">>>DOCTYPE JTL..."}
	respTable := L.NewTable()
	statusint, err := strconv.Atoi(resp["JTLTP-STATUS"])
	if err != nil {
		fmt.Printf("Error converting status to int: %v\n", err)
		return 0
	}

	respTable.RawSetString("JTLTP-STATUS", lua.LNumber(statusint))
	respTable.RawSetString("JTLTP-TYPE", lua.LString(resp["JTLTP-TYPE"]))
	respTable.RawSetString("JTLTP-MSG", lua.LString(resp["JTLTP"]))
	L.Push(respTable)
	return 1
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
	Site = ""
	// Execute script after objects are created and stored
	if err := luaState.DoString(combinedScript); err != nil {
		fmt.Printf("Initial script execution error: %v\n", err)
	}

	return newLocker(objects), objects
}

// AddStyle adds a style to an element
func addStyle(L *lua.LState) int {
	selector := L.ToString(1)
	style := L.ToString(2)

	if element := getElement(selector); element != nil {
		if baseEl, ok := element.(interface{ GetBaseElement() *BaseElement }); ok {
			baseEl.GetBaseElement().AddStyle(style)
		}
	}
	return 0
}

// RemoveAllStyle removes all styles from an element
func removeAllStyle(L *lua.LState) int {
	selector := L.ToString(1)

	if element := getElement(selector); element != nil {
		if baseEl, ok := element.(interface{ GetBaseElement() *BaseElement }); ok {
			baseEl.GetBaseElement().RemoveAllStyle()
		}
	}
	return 0
}

// Add this new function
func executeFrameHandler() {
	if frameHandler != "" && luaState != nil {
		if err := luaState.DoString(frameHandler); err != nil {
			fmt.Printf("Error executing frame handler: %v\n", err)
		}
	}
}

// Add this new function
func executeRequestedFrameHandler() {
	if requestedFrameHandler != "" && luaState != nil {
		if err := luaState.DoString(requestedFrameHandler); err != nil {
			fmt.Printf("Error executing frame handler: %v\n", err)
		}
		requestedFrameHandler = "" // Clear after execution
	}
}

// Add this new function
func setFrameHandler(L *lua.LState) int {
	handler := L.ToString(1)
	frameHandler = handler
	return 0
}

// Add this new function
func setRequestFrameHandler(L *lua.LState) int {
	handler := L.ToString(1)
	requestedFrameHandler = handler
	return 0
}

// Modify drawRenderingState to include frame handler execution
func drawRenderingState(objects []CanvasObject) {
	ObjectsMutex.Lock()
	localObjects := make([]CanvasObject, len(objects))
	copy(localObjects, objects)
	ObjectsMutex.Unlock()

	if len(localObjects) == 0 {
		fmt.Println("No objects to draw")
		return
	}

	// Execute frame handler before drawing objects
	executeFrameHandler()

	// Execute requested frame handler before drawing objects
	executeRequestedFrameHandler()

	for _, obj := range localObjects {
		if obj == nil {
			continue
		}
		obj.Draw()
		obj.CheckClick()
	}
}
