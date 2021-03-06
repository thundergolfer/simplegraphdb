package simplegraphdb

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/thundergolfer/turtle" // TODO: use fork until my PR is merged to handle RDF 'Bags'
)

func (db tripleDb) toString() string {
	return toJSON(db)
}

func toJSON(db interface{}) string {
	bytes, err := json.Marshal(db)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return string(bytes)
}

type tripleDb struct {
	Triples []Entry
}

// Entry is a type used as an intermediary between
// a plaintext/JSON representation of a triple and the
// ID based presentation loaded into the Hexastore
type Entry struct {
	Subject string
	Prop    string
	Object  string
}

// Dictionary is only exported because it's currently tested
// TODO remove need to export this
type Dictionary struct {
	m       map[int]string
	NextKey int
}

// GetKey provides access to the Dictionary's map, returning the
// int ID for a given string val, or (0, false) if the string is not
// in the dictionary
func (dict Dictionary) GetKey(val string) (key int, ok bool) {
	for k, v := range dict.m {
		if v == val {
			return k, true
		}
	}

	return 0, false
}

// Put adds a new string value into the Dictionary.
// Warning: should not be called with the same string twice,
// because that string will get 2 IDs
func (dict *Dictionary) Put(val string) (key int) {
	// Don't call without checking that value doesn't already exist
	dict.m[dict.NextKey] = val
	key = dict.NextKey
	dict.NextKey++
	return
}

// Get accesses the Dictionary's map and returns the result of
// access the given key in that map
// ie. basically a middleman method
func (dict Dictionary) Get(key int) (val string, ok bool) {
	val, ok = dict.m[key]
	return
}

// EntityDict is the Dictionary type for Subjects and Objects
type EntityDict struct {
	Dictionary
}

// PropDict is the Dictionary type for Properties
type PropDict struct {
	Dictionary
}

// NewEntityDict creates and initialises a new EntityDict
func NewEntityDict() *EntityDict {
	var eD EntityDict
	eD.m = make(map[int]string)
	eD.NextKey = 0
	return &eD
}

// NewPropDict creates and initialises and new PropDict
func NewPropDict() *PropDict {
	var pD PropDict
	pD.m = make(map[int]string)
	pD.NextKey = 0
	return &pD
}

// Hexastore is the interface with which the simplesparql query
// engine can interact with the graph data
type Hexastore interface {
	QueryXXX() *[]Triple
	QuerySXX(subjID int) *[]Triple
	QueryXPX(propID int) *[]Triple
	QueryXXO(objID int) *[]Triple
	QuerySPX(subjID, propID int) *[]Triple
	QuerySXO(subjID, objID int) *[]Triple
	QueryXPO(propID, objID int) *[]Triple
	QuerySPO(subjID, propID, objID int) *[]Triple
	Add(subject, property, object, value string)
	GetPropKey(val string) (key int, ok bool)
	GetEntityKey(val string) (key int, ok bool)
	ResolveEntity(id int) string
	ResolveProp(id int) string
}

// HexastoreDB is the triple-store data structure
// driving the database
type HexastoreDB struct {
	SPO      map[int]map[int]map[int]string
	SOP      map[int]map[int]map[int]string
	PSO      map[int]map[int]map[int]string
	POS      map[int]map[int]map[int]string
	OSP      map[int]map[int]map[int]string
	OPS      map[int]map[int]map[int]string
	entities *EntityDict
	props    *PropDict
}

func newHexastore() *HexastoreDB {
	var store HexastoreDB
	store.SPO = make(map[int]map[int]map[int]string)
	store.SOP = make(map[int]map[int]map[int]string)
	store.PSO = make(map[int]map[int]map[int]string)
	store.POS = make(map[int]map[int]map[int]string)
	store.OSP = make(map[int]map[int]map[int]string)
	store.OPS = make(map[int]map[int]map[int]string)
	store.entities = NewEntityDict()
	store.props = NewPropDict()

	return &store
}

// Triple is the data structure used by Hexastore to
// store a triple
type Triple struct {
	Subject int
	Prop    int
	Object  int
	Value   string
}

// MakeTriple is a cosmetic only Triple initialiser method
func MakeTriple(subject, prop, object int, value string) *Triple {
	t := Triple{Subject: subject, Prop: prop, Object: object, Value: value}
	return &t
}

// PresentableResults converts a slice of Triple objects to a slice of strings with
// triple component IDs converted to their string value
func PresentableResults(results *[]Triple, hexastore *HexastoreDB) *[]string {
	presentables := []string{}
	for _, triple := range *results {
		presentables = append(presentables, PresentTriple(&triple, hexastore.props, hexastore.entities))
	}

	return &presentables
}

// PresentTriple is essentially triple.ToString(), but we use this because the store's
// two dictionaries are needed
// TODO refactor to use store.Resolve*() methods
func PresentTriple(t *Triple, props *PropDict, entities *EntityDict) string {
	return fmt.Sprintf("%s -> %s -> %s", entities.m[t.Subject], props.m[t.Prop], entities.m[t.Object])
}

// GetPropKey finds the integer id for a particular string value which is a property
func (store *HexastoreDB) GetPropKey(val string) (int, bool) {
	return store.props.GetKey(val)
}

// GetEntityKey finds the integer id for a particular string value which is an entity (object or subject)
func (store *HexastoreDB) GetEntityKey(val string) (int, bool) {
	return store.entities.GetKey(val)
}

// ResolveEntity finds the string value for a given entity ID
func (store *HexastoreDB) ResolveEntity(id int) string {
	val, _ := store.entities.Get(id)
	return val
}

// ResolveProp finds the string value for a given property ID
func (store *HexastoreDB) ResolveProp(id int) string {
	val, _ := store.props.Get(id)
	return val
}

// Add introduces a new triple into the hexastore database
func (store *HexastoreDB) Add(subject, property, object, value string) {
	subjectID, propID, objectID := store.MapStringsToIds(subject, property, object)
	triple := MakeTriple(subjectID, propID, objectID, value)
	store.add(triple)
}

// MapIdsToStrings finds the string values for each component ID of a triple
func (store *HexastoreDB) MapIdsToStrings(subjID, propID, objectID int) (string, string, string) {
	subject, _ := store.entities.Get(subjID)
	object, _ := store.entities.Get(objectID)
	prop, _ := store.props.Get(propID)

	return subject, prop, object
}

// MapStringsToIds finds the IDs for each string value of a triple, and creates and entry and returns
// the new ID if a string value doesn't already exist
func (store *HexastoreDB) MapStringsToIds(subject, property, object string) (int, int, int) {
	var ok bool
	var objID, subjID, propID int
	subjID, ok = store.entities.GetKey(subject)
	if !ok {
		subjID = store.entities.Put(subject)
	}
	propID, ok = store.props.GetKey(property)
	if !ok {
		propID = store.props.Put(property)
	}
	objID, ok = store.entities.GetKey(object)
	if !ok {
		objID = store.entities.Put(object)
	}

	return subjID, propID, objID
}

func (store *HexastoreDB) add(t *Triple) {
	var s, p, o int = t.Subject, t.Prop, t.Object
	v := t.Value

	if store.SPO[s] == nil {
		store.SPO[s] = make(map[int]map[int]string)
	}
	if store.SPO[s][p] == nil {
		store.SPO[s][p] = make(map[int]string)
	}
	store.SPO[s][p][o] = v

	if store.SOP[s] == nil {
		store.SOP[s] = make(map[int]map[int]string)
	}
	if store.SOP[s][o] == nil {
		store.SOP[s][o] = make(map[int]string)
	}
	store.SOP[s][o][p] = v

	if store.PSO[p] == nil {
		store.PSO[p] = make(map[int]map[int]string)
	}
	if store.PSO[p][s] == nil {
		store.PSO[p][s] = store.SPO[s][p]
	}
	if store.POS[p] == nil {
		store.POS[p] = make(map[int]map[int]string)
	}
	if store.POS[p][o] == nil {
		store.POS[p][o] = make(map[int]string)
	}
	store.POS[p][o][s] = v

	if store.OSP[o] == nil {
		store.OSP[o] = make(map[int]map[int]string)
	}
	if store.OSP[o][s] == nil {
		store.OSP[o][s] = store.SOP[s][o]
	}
	if store.OPS[o] == nil {
		store.OPS[o] = make(map[int]map[int]string)
	}
	if store.OPS[o][p] == nil {
		store.OPS[o][p] = store.POS[p][o]
	}
}

func (store *HexastoreDB) remove(t *Triple) {
	var s, p, o int = t.Subject, t.Prop, t.Object
	var subjectIndx = store.SPO
	var pred = subjectIndx[s]

	if pred != nil {
		var obj = pred[p]
		if obj != nil {
			if _, ok := obj[o]; ok {
				delete(store.SPO[s][p], o)
				delete(store.SOP[s][o], p)
				delete(store.PSO[p][s], o)
				delete(store.POS[p][o], s)
				delete(store.OSP[o][s], p)
				delete(store.OPS[o][p], s)
			}
		}
	}
}

/*
 * Query methods
 */

// QuerySXX allows for querying the hexastore specifying only a Subject entity ID
func (store *HexastoreDB) QuerySXX(subjID int) *[]Triple {
	res := []Triple{}
	relevant := store.SPO[subjID]

	if relevant == nil {
		return &[]Triple{}
	}

	for prop, objMap := range relevant {
		for obj, value := range objMap {
			currTriple := MakeTriple(subjID, prop, obj, value)
			res = append(res, *currTriple)
		}
	}

	return &res
}

// QuerySPX allows for querying the hexastore specifying only a Subject entity ID
// and a Property ID
func (store *HexastoreDB) QuerySPX(subjID, propID int) *[]Triple {
	res := []Triple{}
	relevant := store.SPO[subjID]

	if relevant == nil {
		return &[]Triple{}
	}

	properties := relevant[propID]

	for objID, value := range properties {
		currTriple := MakeTriple(subjID, propID, objID, value)
		res = append(res, *currTriple)
	}

	return &res
}

// QuerySXO allows for querying the hexastore specifying only a Subject entity ID
// and an Object entity ID
func (store *HexastoreDB) QuerySXO(subjID, objID int) *[]Triple {
	return &[]Triple{}
}

// QueryXPX allows for querying the hexastore specifying only a Property ID
func (store *HexastoreDB) QueryXPX(propID int) *[]Triple {
	res := []Triple{}
	relevant := store.PSO[propID]

	if relevant == nil {
		return &[]Triple{}
	}

	for subjID, objMap := range relevant {
		for objID, value := range objMap {
			currTriple := MakeTriple(subjID, propID, objID, value)
			res = append(res, *currTriple)
		}
	}

	return &res
}

// QueryXPO allows for querying the hexastore specifying only a Property ID
// and an Object entity ID
func (store *HexastoreDB) QueryXPO(propID, objID int) *[]Triple {
	res := []Triple{}
	relevant := store.POS[propID]

	if relevant == nil {
		return &[]Triple{}
	}

	subjects := relevant[objID]

	if subjects == nil {
		return &[]Triple{}
	}

	for subjID, value := range subjects {
		currTriple := MakeTriple(subjID, propID, objID, value)
		res = append(res, *currTriple)
	}

	return &res
}

// QueryXXO allows for querying the hexastore specifying only an Object entity ID
func (store *HexastoreDB) QueryXXO(objID int) *[]Triple {
	res := []Triple{}
	relevant := store.OPS[objID]

	if relevant == nil {
		return &[]Triple{}
	}

	for propID, subjMap := range relevant {
		for subjID, value := range subjMap {
			currTriple := MakeTriple(subjID, propID, objID, value)
			res = append(res, *currTriple)
		}
	}

	return &res
}

// QuerySPO allows for querying the hexastore for a specific triple
func (store *HexastoreDB) QuerySPO(subjID, propID, objID int) *[]Triple {
	if value, ok := store.SPO[subjID][propID][objID]; ok {
		triple := MakeTriple(subjID, propID, objID, value)
		return &[]Triple{*triple}
	}

	return &[]Triple{}
}

// QueryXXX allows for a 'wildcard' query of the hexastore, returning all triples
func (store *HexastoreDB) QueryXXX() *[]Triple {
	res := []Triple{}

	for subjID, propMap := range store.SPO {
		for propID, objMap := range propMap {
			for objID, value := range objMap {
				currTriple := MakeTriple(subjID, propID, objID, value)
				res = append(res, *currTriple)
			}
		}
	}

	return &res
}

/*
 * End Query methods
 */

func loadHexastore(db tripleDb, store *HexastoreDB) error {
	for _, entry := range db.Triples {
		val := "xxxx" // TODO

		store.Add(entry.Subject, entry.Prop, entry.Object, val)
	}

	return nil
}

// InitHexastoreFromJSON creates a new hexastore and fills it with triples
// from a valid JSON file. The schema is:
// {"triples": [
//  	{"subject": <STRING>, "prop": <STRING>, "object": <STRING>},
//    ...
//    ]
// }
func InitHexastoreFromJSON(dbFilePath string) (*HexastoreDB, error) {
	var db tripleDb

	dat, err := ioutil.ReadFile(dbFilePath)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(dat, &db)

	store := newHexastore()
	err = loadHexastore(db, store)
	if err != nil {
		return nil, err
	}

	return store, nil
}

// InitHexastoreFromJSONRows creates a new hexastore and fills it with triples
// from a file with one JSON object (a triple) per line. The schema is:
//
// {"subject": <STRING>, "prop": <STRING>, "object": <STRING>}
// {"subject": <STRING>, "prop": <STRING>, "object": <STRING>}
// {"subject": <STRING>, "prop": <STRING>, "object": <STRING>}
func InitHexastoreFromJSONRows(dbFilePath string) (Hexastore, error) {
	db := tripleDb{Triples: []Entry{}}
	var entry Entry

	// dat, err := ioutil.ReadFile(dbFilePath)
	// if err != nil {
	// 	return nil, err
	// }
	file, err := os.Open(dbFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		json.Unmarshal(scanner.Bytes(), &entry)
		db.Triples = append(db.Triples, entry)
	}

	store := newHexastore()
	err = loadHexastore(db, store)
	if err != nil {
		return nil, err
	}

	return store, nil
}

// InitHexastoreFromTurtle creates a new hexastore and fills it with triples
// from a file in Terse RDF Triple Language, or 'Turtle' (https://www.w3.org/TeamSubmission/turtle/)
func InitHexastoreFromTurtle(dbFilePath string) (Hexastore, error) {
	db := tripleDb{}

	dat, err := ioutil.ReadFile(dbFilePath)
	if err != nil {
		return nil, err
	}

	triples, err := turtle.Parse(dat)
	if err != nil {
		return nil, err
	}

	db.Triples = make([]Entry, len(triples))
	for i, t := range triples {
		entry := Entry{Subject: t.Subj, Prop: t.Pred, Object: t.Obj}
		db.Triples[i] = entry
	}

	store := newHexastore()
	err = loadHexastore(db, store)
	if err != nil {
		return nil, err
	}

	return store, nil
}
