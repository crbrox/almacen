package almacen

import (
	"testing"
)

func makeMongoTest(functest func(s Store, t *testing.T)) func(t *testing.T) {
	return func(t *testing.T) {
		store := &MongoEntityStore{}
		err := store.Start(&Config{MongoURL: "localhost"})
		contextTest.session = store.session
		if err != nil {
			t.Fatal(err)
		}
		contextTest.session.DB("").C(collectionTest).DropCollection()
		defer store.Stop()

		functest(store, t)

	}
}

func TestMongoStoreFindAll(t *testing.T) {
	makeMongoTest(testFindAll)(t)
}

func TestMongoStoreFindByID(t *testing.T) {
	makeMongoTest(testFindByID)(t)
}

func TestMongoStoreFindByIDNotFound(t *testing.T) {
	makeMongoTest(testFindByIDNotFound)(t)
}

func TestMongoStoreSaveIDNotString(t *testing.T) {
	makeMongoTest(testSaveIDNotString)(t)
}

func TestMongoStoreDelete(t *testing.T) {
	makeMongoTest(testDelete)(t)
}

func TestMongoStoreFindField(t *testing.T) {
	makeMongoTest(testFindField)(t)
}

func TestMongoStoreFindFieldNested(t *testing.T) {
	makeMongoTest(testFindFieldNested)(t)
}

func TestMongoStoreUpdateField(t *testing.T) {
	makeMongoTest(testUpdateField)(t)
}

func TestMongoStoreUpdateFieldTraversingErr(t *testing.T) {
	makeMongoTest(testUpdateFieldTraversingErr)(t)
}

func TestMongoStoreDeleteField(t *testing.T) {
	makeMongoTest(testDeleteField)(t)
}

func TestMongoStoreDeleteFieldTraversingErr(t *testing.T) {
	makeMongoTest(testDeleteFieldTraversingErr)(t)
}

func TestMongoStoreUpdateFieldRoot(t *testing.T) {
	makeMongoTest(testUpdateFieldRoot)(t)
}

func TestMongoStoreDeleteFieldRoot(t *testing.T) {
	makeMongoTest(testDeleteFieldRoot)(t)
}

func TestStartErr(t *testing.T) {
	store := &MongoEntityStore{}
	err := store.Start(&Config{MongoURL: "piticlin?a_very_rare_option=0"})
	if err == nil {
		t.Errorf("start error: wanted <something>, got nil")
	}
}