package simplegraphdb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (db Db) toString() string {
	return toJson(db)
}

func toJson(db interface{}) string {
	bytes, err := json.Marshal(db)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return string(bytes)
}

type Db struct {
	Triples []Entry
}

type Entry struct {
	Subject string
	Prop    string
	Object  string
}

type Dictionary struct {
	m       map[int]string
	NextKey int
}

func (dict Dictionary) GetKey(val string) (key int, ok bool) {
	for k, v := range dict.m {
		if v == val {
			return k, true
		}
	}

	return 0, false
}

func (dict *Dictionary) Put(val string) (key int) {
	// Don't call without checking that value doesn't already exist
	dict.m[dict.NextKey] = val
	key = dict.NextKey
	dict.NextKey += 1
	return
}

func (dict Dictionary) Get(key int) (val string, ok bool) {
	val, ok = dict.m[key]
	return
}

type EntityDict struct {
	Dictionary
}

type PropDict struct {
	Dictionary
}

func NewEntityDict() *EntityDict {
	var eD EntityDict
	eD.m = make(map[int]string)
	eD.NextKey = 0
	return &eD
}

func NewPropDict() *PropDict {
	var pD PropDict
	pD.m = make(map[int]string)
	pD.NextKey = 0
	return &pD
}

type Hexastore struct {
	SPO      map[int]map[int]map[int]string
	SOP      map[int]map[int]map[int]string
	PSO      map[int]map[int]map[int]string
	POS      map[int]map[int]map[int]string
	OSP      map[int]map[int]map[int]string
	OPS      map[int]map[int]map[int]string
	entities *EntityDict
	props    *PropDict
}

func newHexastore() *Hexastore {
	var store Hexastore
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

type Triple struct {
	Subject int
	Prop    int
	Object  int
	Value   string
}

func MakeTriple(subject, prop, object int, value string) *Triple {
	t := Triple{Subject: subject, Prop: prop, Object: object, Value: value}
	return &t
}

func PresentableResults(results *[]Triple, hexastore *Hexastore) *[]string {
	presentables := []string{}
	for _, triple := range *results {
		presentables = append(presentables, PresentTriple(&triple, hexastore.props, hexastore.entities))
	}

	return &presentables
}

func PresentTriple(t *Triple, props *PropDict, entities *EntityDict) string {
	return fmt.Sprintf("%s -> %s -> %s", entities.m[t.Subject], props.m[t.Prop], entities.m[t.Object])
}

func (store Hexastore) ResolveEntity(id int) string {
	val, _ := store.entities.Get(id)
	return val
}

func (store Hexastore) ResolveProp(id int) string {
	val, _ := store.props.Get(id)
	return val
}

func (store Hexastore) Add(subject, property, object, value string) bool {
	subjectId, propId, objectId := store.MapStringsToIds(subject, property, object)
	triple := MakeTriple(subjectId, propId, objectId, value)

	store.add(triple)

	return true
}

func (store Hexastore) MapIdsToStrings(subjId, propId, objectId int) (string, string, string) {
	subject, _ := store.entities.Get(subjId)
	object, _ := store.entities.Get(objectId)
	prop, _ := store.props.Get(propId)

	return subject, prop, object
}

func (store Hexastore) MapStringsToIds(subject, property, object string) (subjId, propId, objId int) {
	var ok bool
	subjId, ok = store.entities.GetKey(subject)
	if !ok {
		subjId = store.entities.Put(subject)
	}
	propId, ok = store.props.GetKey(property)
	if !ok {
		propId = store.props.Put(property)
	}
	objId, ok = store.entities.GetKey(object)
	if !ok {
		objId = store.entities.Put(object)
	}

	_ = ok // TODO fix this weirdness

	return
}

func (store Hexastore) add(t *Triple) {
	var s, p, o int = t.Subject, t.Prop, t.Object
	var v string = t.Value

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

func (store Hexastore) remove(t *Triple) {
	var s, p, o int = t.Subject, t.Prop, t.Object

	var subject_indx = store.SPO

	var pred = subject_indx[s]

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

func (store Hexastore) QuerySXX(subjId int) *[]Triple {
	res := []Triple{}
	relevant := store.SPO[subjId]

	if relevant == nil {
		return &[]Triple{}
	}

	for prop, objMap := range relevant {
		for obj, value := range objMap {
			currTriple := MakeTriple(subjId, prop, obj, value)
			res = append(res, *currTriple)
		}
	}

	return &res
}

func (store Hexastore) QuerySPX(subjId, propId int) *[]Triple {
	res := []Triple{}
	relevant := store.SPO[subjId]

	if relevant == nil {
		return &[]Triple{}
	}

	properties := relevant[propId]

	for objId, value := range properties {
		currTriple := MakeTriple(subjId, propId, objId, value)
		res = append(res, *currTriple)
	}

	return &res
}

func (store Hexastore) QuerySXO(subjId, objId int) *[]Triple {
	return &[]Triple{}
}

func (store Hexastore) QueryXPX(propId int) *[]Triple {
	res := []Triple{}
	relevant := store.PSO[propId]

	if relevant == nil {
		return &[]Triple{}
	}

	for subjId, objMap := range relevant {
		for objId, value := range objMap {
			currTriple := MakeTriple(subjId, propId, objId, value)
			res = append(res, *currTriple)
		}
	}

	return &res
}

func (store Hexastore) QueryXPO(propId, objId int) *[]Triple {
	res := []Triple{}
	relevant := store.POS[propId]

	if relevant == nil {
		return &[]Triple{}
	}

	subjects := relevant[objId]

	if subjects == nil {
		return &[]Triple{}
	}

	for subjId, value := range subjects {
		currTriple := MakeTriple(subjId, propId, objId, value)
		res = append(res, *currTriple)
	}

	return &res
}

func (store Hexastore) QueryXXO(objId int) *[]Triple {
	res := []Triple{}
	relevant := store.OPS[objId]

	if relevant == nil {
		return &[]Triple{}
	}

	for propId, subjMap := range relevant {
		for subjId, value := range subjMap {
			currTriple := MakeTriple(subjId, propId, objId, value)
			res = append(res, *currTriple)
		}
	}

	return &res
}

func (store Hexastore) QuerySPO(subjId, propId, objId int) *[]Triple {
	if value, ok := store.SPO[subjId][propId][objId]; ok {
		triple := MakeTriple(subjId, propId, objId, value)
		return &[]Triple{*triple}
	}

	return &[]Triple{}
}

func (store Hexastore) QueryXXX() *[]Triple {
	res := []Triple{}

	for subjId, propMap := range store.SPO {
		for propId, objMap := range propMap {
			for objId, value := range objMap {
				currTriple := MakeTriple(subjId, propId, objId, value)
				res = append(res, *currTriple)
			}
		}
	}

	return &res
}

func (store Hexastore) dump() *[]Triple {
	return store.QueryXXX()
}

func loadHexastore(db Db, store *Hexastore) error {
	for _, entry := range db.Triples {
		val := "xxxx" // TODO

		store.Add(entry.Subject, entry.Prop, entry.Object, val)
	}

	return nil
}

func InitTestHexastore() *Hexastore {
	return InitHexastoreFromJson("./db.json")
}

func InitHexastoreFromJson(dbFilePath string) *Hexastore {
	var db Db

	dat, err := ioutil.ReadFile(dbFilePath)
	check(err)

	json.Unmarshal(dat, &db)

	store := newHexastore()
	_ = loadHexastore(db, store)

	return store
}
