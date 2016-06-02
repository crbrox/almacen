package almacen

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var store Store

func SetStore(es Store) {
	store = es
}

type MongoConfig struct {
	URL string
}

type MongoEntityStore struct {
	session *mgo.Session
}

type Store interface {
	FindAll(ctx *context, collection string) ([]map[string]interface{}, error)
	FindByID(ctx *context, collection, id string) (map[string]interface{}, error)
	Save(ctx *context, collection string, ent map[string]interface{}) error
	Delete(ctx *context, collection, id string) error
	FindField(ctx *context, collection, id, field string) (interface{}, error)
	UpdateField(ctx *context, collection, id, field string, value interface{}) error
	DeleteField(ctx *context, collection, id, field string) error
}

func (mes *MongoEntityStore) Start(config *Config) (err error) {
	// AÑADIR Check contra inicializacion doble/concurrente
	mes.session, err = mgo.Dial(config.MongoURL)
	if err != nil {
		return &Error{statusCode: 500, message: "MongoEntityStore.Start: " + err.Error()}
	}

	// Optional. Switch the session to a monotonic behavior.
	mes.session.SetMode(mgo.Monotonic, true)

	// Verificacion de indices ...
	// Opcional ...

	return nil
}

func (mes *MongoEntityStore) Stop() {
	mes.session.Close()
}

func (*MongoEntityStore) FindAll(ctx *context, collection string) ([]map[string]interface{}, error) {
	var list []map[string]interface{}
	err := ctx.session.DB("").C(collection).Find(nil).All(&list)
	return list, err
}

func (*MongoEntityStore) FindByID(ctx *context, collection, id string) (map[string]interface{}, error) {
	var list []map[string]interface{}
	err := ctx.session.DB("").C(collection).FindId(id).All(&list)
	if len(list) == 0 && err == nil {
		return nil, ErrNotFound
	}
	return list[0], err
}

func (*MongoEntityStore) Save(ctx *context, collection string, ent map[string]interface{}) error {
	id, isString := ent["_id"].(string)
	if !isString {
		return ErrIdNotString
	}
	_, err := ctx.session.DB("").C(collection).UpsertId(id, ent)
	return err
}

func (*MongoEntityStore) Delete(ctx *context, collection, id string) error {
	err := ctx.session.DB("").C(collection).Remove(bson.M{"_id": id})
	return err
}

func (*MongoEntityStore) FindField(ctx *context, collection, id, field string) (interface{}, error) {
	var result map[string]interface{}
	const resultKey = "V"
	/*
		err := ctx.session.DB("").C(collection).
			Find(bson.M{"_id": id}).
			Select((bson.M{field: 1, "_id": 0})).
			One(&entity)

		if err != nil {
			return nil, err
		}
		val := extractField(field, entity)
	*/
	err := ctx.session.DB("").C(collection).
		Pipe([]bson.M{
			{"$match": bson.M{"_id": id}},
			{"$project": bson.M{resultKey: "$" + field, "_id": false}},
		}).One(&result)
	if err != nil {
		return nil, err
	}
	val, present := result[resultKey]
	if !present {
		return nil, ErrNotFound
	}
	return val, nil
}

func (*MongoEntityStore) UpdateField(ctx *context, collection, id, field string, value interface{}) error {
	err := ctx.session.DB("").C(collection).Update(
		bson.M{"_id": id},
		bson.M{"$set": bson.M{field: value}})
	if err, ok := err.(*mgo.LastError); ok {
		if err.Code == 16837 { // Wrong traverse
			return ErrTraversingObject
		}
	}
	return err
}

func (*MongoEntityStore) DeleteField(ctx *context, collection, id, field string) error {
	err := ctx.session.DB("").C(collection).Update(
		bson.M{"_id": id},
		bson.M{"$unset": bson.M{field: 1}})
	return err
}

/*
func extractField(field string, o map[string]interface{}) interface{} {
	var (
		current map[string]interface{}
		value   interface{}
	)
	fs := strings.Split(field, ".")
	current = o
	for _, f := range fs[:len(fs)-1] {
		value = current[f]
		switch value := value.(type) {
		case map[string]interface{}:
			current = value
		default:
			return nil
		}
	}
	return current[fs[len(fs)-1]]
}
*/
