package almacen

import (
	"reflect"
	"sort"
	"testing"
)

const collectionTest = "collectionTest"

type sortableEntitySlice []map[string]interface{}

func (s sortableEntitySlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s sortableEntitySlice) Len() int      { return len(s) }
func (s sortableEntitySlice) Less(i, j int) bool {
	idI := s[i]["_id"].(string)
	idJ := s[j]["_id"].(string)
	return idI < idJ
}

var entitiesTest []map[string]interface{}
var contextTest = &context{}

func initTest() {
	entitiesTest = []map[string]interface{}{
		{"_id": "e1", "temperature": 12.34},
		{"_id": "e2", "temperature": 56.78, "label": "very hot"},
		{"_id": "e3", "location": map[string]interface{}{"lon": -3.7025600, "lat": 40.4165000}},
		{"_id": "e4", "nested": map[string]interface{}{"sub1": map[string]interface{}{"sub2": map[string]interface{}{"sub3": map[string]interface{}{}}}}},
	}
}

func populateTest(s Store, t *testing.T) {
	initTest()
	for _, e := range entitiesTest {
		err := s.Save(contextTest, collectionTest, e)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func testFindAll(s Store, t *testing.T) {
	populateTest(s, t)
	res, err := s.FindAll(contextTest, collectionTest)
	if err != nil {
		t.Fatal(err)
	}
	sort.Sort(sortableEntitySlice(res))
	sort.Sort(sortableEntitySlice(entitiesTest))
	if !reflect.DeepEqual(res, entitiesTest) {
		t.Errorf("expected %v, got %v", entitiesTest, res)
	}
}

func testFindByID(s Store, t *testing.T) {
	populateTest(s, t)
	res, err := s.FindByID(contextTest, collectionTest, entitiesTest[0]["_id"].(string))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(res, entitiesTest[0]) {
		t.Errorf("expected %v, got %v", entitiesTest, res)
	}
}

func testFindByIDNotFound(s Store, t *testing.T) {
	populateTest(s, t)
	_, err := s.FindByID(contextTest, collectionTest, "shouldnotexist")
	if err != ErrNotFound {
		t.Errorf("expected %v, got %v", ErrNotFound, err)
	}
}

func testSaveIDNotString(s Store, t *testing.T) {
	initTest()
	entitiesTest[0]["_id"] = 42
	err := s.Save(contextTest, collectionTest, entitiesTest[0])
	if err != ErrIdNotString {
		t.Errorf("expected %v, got %v", ErrIdNotString, err)
	}
}

func testDelete(s Store, t *testing.T) {
	populateTest(s, t)
	if err := s.Delete(contextTest, collectionTest, entitiesTest[1]["_id"].(string)); err != nil {
		t.Fatal(err)
	}
	res, err := s.FindAll(contextTest, collectionTest)
	if err != nil {
		t.Fatal(err)
	}
	entitiesTest = append(entitiesTest[0:1], entitiesTest[2:]...)
	sort.Sort(sortableEntitySlice(res))
	sort.Sort(sortableEntitySlice(entitiesTest))
	if !reflect.DeepEqual(res, entitiesTest) {
		t.Errorf("expected %v, got %v", entitiesTest, res)
	}
}

func testFindField(s Store, t *testing.T) {
	populateTest(s, t)
	for _, entity := range entitiesTest {
		key := entity["_id"].(string)
		for field, value := range entity {
			if field != "_id" {
				res, err := s.FindField(contextTest, collectionTest, key, field)
				if err != nil {
					t.Errorf("%#v value: %#v, id: %#v, field: %#v", err, value, key, field)
				}

				if !reflect.DeepEqual(res, value) {
					t.Errorf("expected %v, got %v", value, res)
				}
			}
		}
	}
}

func testFindFieldNested(s Store, t *testing.T) {
	id := "ID"
	elZ := map[string]interface{}{"z": 12}
	elY := map[string]interface{}{"y": elZ}
	el := map[string]interface{}{"_id": "ID", "x": elY}
	cases := map[string]interface{}{
		"x":     elY,
		"x.y":   elZ,
		"x.y.z": 12,
	}

	if err := s.Save(contextTest, collectionTest, el); err != nil {
		t.Fatal(err)
	}

	for field, value := range cases {
		res, err := s.FindField(contextTest, collectionTest, id, field)
		if err != nil {
			t.Errorf("%#v value: %#v, id: %#v, field: %#v", err, value, id, field)
		}
		if !reflect.DeepEqual(res, value) {
			t.Errorf("expected %v, got %v", value, res)
		}
	}

	// Not found
	cases["not.exist.i.hope"] = nil
	for field := range cases {
		field := field + "__xx"
		res, err := s.FindField(contextTest, collectionTest, id, field)
		if err != ErrNotFound {
			t.Errorf("expected %v, got %v", ErrNotFound, res)
		}
	}
}

func testUpdateField(s Store, t *testing.T) {
	original := map[string]interface{}{"_id": "ID", "x": map[string]interface{}{"y": map[string]interface{}{"z": 12}}}
	expected := map[string]interface{}{"_id": "ID", "x": map[string]interface{}{"y": map[string]interface{}{"z": "CHANGED"}}}

	if err := s.Save(contextTest, collectionTest, original); err != nil {
		t.Fatal(err)
	}
	if err := s.UpdateField(contextTest, collectionTest, "ID", "x.y.z", "CHANGED"); err != nil {
		t.Fatal(err)
	}
	res, err := s.FindByID(contextTest, collectionTest, "ID")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(res, expected) {
		t.Errorf("expected %v, got %v", expected, res)
	}
}

func testUpdateFieldTraversingErr(s Store, t *testing.T) {
	original := map[string]interface{}{"_id": "ID", "x": map[string]interface{}{"y": map[string]interface{}{"z": 12}}}

	if err := s.Save(contextTest, collectionTest, original); err != nil {
		t.Fatal(err)
	}
	err := s.UpdateField(contextTest, collectionTest, "ID", "x.y.z.t", "CHANGED")
	expected := ErrTraversingObject
	if !reflect.DeepEqual(err, expected) {
		t.Errorf("expected %#v, got %#v", expected, err)
	}
}

func testDeleteField(s Store, t *testing.T) {
	original := map[string]interface{}{"_id": "ID", "x": map[string]interface{}{"y": map[string]interface{}{"z": 12}}, "a": "A"}
	expected := map[string]interface{}{"_id": "ID", "x": map[string]interface{}{"y": map[string]interface{}{}}, "a": "A"}

	if err := s.Save(contextTest, collectionTest, original); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteField(contextTest, collectionTest, "ID", "x.y.z"); err != nil {
		t.Fatal(err)
	}
	res, err := s.FindByID(contextTest, collectionTest, "ID")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(res, expected) {
		t.Errorf("expected %v, got %v", expected, res)
	}
}

func testDeleteFieldTraversingErr(s Store, t *testing.T) {
	original := map[string]interface{}{"_id": "ID", "x": map[string]interface{}{"y": map[string]interface{}{"z": 12}}}

	if err := s.Save(contextTest, collectionTest, original); err != nil {
		t.Fatal(err)
	}
	err := s.DeleteField(contextTest, collectionTest, "ID", "x.y.z.t")
	if !reflect.DeepEqual(err, nil) {
		t.Errorf("expected %v, got %v", nil, err)
	}
}

func testUpdateFieldRoot(s Store, t *testing.T) {
	original := map[string]interface{}{"_id": "ID", "x": "12", "y": 12}
	expected := map[string]interface{}{"_id": "ID", "x": 21, "y": 12}

	if err := s.Save(contextTest, collectionTest, original); err != nil {
		t.Fatal(err)
	}
	if err := s.UpdateField(contextTest, collectionTest, "ID", "x", 21); err != nil {
		t.Fatal(err)
	}
	res, err := s.FindByID(contextTest, collectionTest, "ID")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(res, expected) {
		t.Errorf("expected %v, got %v", expected, res)
	}
}

func testDeleteFieldRoot(s Store, t *testing.T) {
	original := map[string]interface{}{"_id": "ID", "x": "X", "a": "A"}
	expected := map[string]interface{}{"_id": "ID", "a": "A"}

	if err := s.Save(contextTest, collectionTest, original); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteField(contextTest, collectionTest, "ID", "x"); err != nil {
		t.Fatal(err)
	}
	res, err := s.FindByID(contextTest, collectionTest, "ID")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(res, expected) {
		t.Errorf("expected %v, got %v", expected, res)
	}
}
