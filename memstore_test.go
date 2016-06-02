package almacen

import (
	"testing"
)

func TestMemStoreFindAll(t *testing.T) {
	m := NewMemStore()
	testFindAll(m, t)
}

func TestMemStoreFindByID(t *testing.T) {
	m := NewMemStore()
	testFindByID(m, t)
}

func TestMemStoreFindByIDNotFound(t *testing.T) {
	m := NewMemStore()
	testFindByIDNotFound(m, t)
}

func TestMemStoreSaveIDNotString(t *testing.T) {
	m := NewMemStore()
	testSaveIDNotString(m, t)
}

func TestMemStoreDelete(t *testing.T) {
	m := NewMemStore()
	testDelete(m, t)
}

func TestMemStoreFindField(t *testing.T) {
	m := NewMemStore()
	testFindField(m, t)
}
func TestMemStoreFindFieldNested(t *testing.T) {
	m := NewMemStore()
	testFindFieldNested(m, t)
}

func TestMemStoreUpdateField(t *testing.T) {
	m := NewMemStore()
	testUpdateField(m, t)
}

func TestMemStoreUpdateFieldTraversingErr(t *testing.T) {
	m := NewMemStore()
	testUpdateFieldTraversingErr(m, t)
}

func TestMemStoreDeleteField(t *testing.T) {
	m := NewMemStore()
	testDeleteField(m, t)
}

func TestMemStoreDeleteFieldTraversingErr(t *testing.T) {
	m := NewMemStore()
	testDeleteFieldTraversingErr(m, t)
}

func TestMemStoreUpdateFieldRoot(t *testing.T) {
	m := NewMemStore()
	testUpdateFieldRoot(m, t)
}

func TestMemStoreDeleteFieldRoot(t *testing.T) {
	m := NewMemStore()
	testDeleteFieldRoot(m, t)
}
