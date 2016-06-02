package almacen

import (
	"strings"
	"sync"
)

type MemStore struct {
	db map[string]map[string]map[string]interface{}
	mu sync.RWMutex
}

func NewMemStore() *MemStore {
	return &MemStore{db: make(map[string]map[string]map[string]interface{})}
}
func (ms *MemStore) getCol(collection string) map[string]map[string]interface{} {
	col := ms.db[collection]
	if col == nil {
		col = make(map[string]map[string]interface{})
		ms.db[collection] = col
	}
	return col
}

func (ms *MemStore) FindAll(ctx *context, collection string) ([]map[string]interface{}, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	var list []map[string]interface{}
	for _, e := range ms.getCol(collection) {
		list = append(list, e)
	}
	return list, nil
}

func (ms *MemStore) FindByID(ctx *context, collection, id string) (map[string]interface{}, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	var obj, found = ms.getCol(collection)[id]
	if !found {
		return nil, ErrNotFound
	}
	return obj, nil
}

func (ms *MemStore) Save(ctx *context, collection string, ent map[string]interface{}) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	key, isString := ent["_id"].(string)
	if !isString {
		return ErrIdNotString
	}
	ms.getCol(collection)[key] = ent
	return nil
}

func (ms *MemStore) Delete(ctx *context, collection, id string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	delete(ms.getCol(collection), id)
	return nil
}

func (ms *MemStore) FindField(ctx *context, collection, id, field string) (interface{}, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	element, father := ms.traverse(collection, id, field)
	if father == nil {
		return nil, ErrNotFound
	}
	value, isPresent := father[element]
	if !isPresent {
		return nil, ErrNotFound
	}
	return value, nil
}

func (ms *MemStore) UpdateField(ctx *context, collection, id, fields string, value interface{}) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	field, father := ms.traverse(collection, id, fields)
	if father == nil {
		return ErrTraversingObject
	}
	father[field] = value
	return nil
}

func (ms *MemStore) DeleteField(ctx *context, collection, id, fields string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	field, father := ms.traverse(collection, id, fields)
	if father != nil {
		delete(father, field)
	}
	return nil
}

func (ms *MemStore) traverse(collection, id, fields string) (element string, father map[string]interface{}) {
	col := ms.getCol(collection)
	root := col[id]
	fSlice := strings.Split(fields, ".")
	element = fSlice[len(fSlice)-1] // last element in x.y.z -> z
	fSlice = fSlice[:len(fSlice)-1] // path to z, [x,y]
	if len(fSlice) == 0 {           // first level field, root is the father
		return element, root
	}

	father = traverseAux(root, fSlice) // find the father of element
	return element, father
}

func traverseAux(m map[string]interface{}, fields []string) (father map[string]interface{}) {
	f, isObject := m[fields[0]].(map[string]interface{})
	if !isObject || f == nil {
		return nil
	}
	if len(fields) == 1 {
		return f
	}
	return traverseAux(f, fields[1:])
}
