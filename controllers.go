package almacen

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func AddRoutes(router *httprouter.Router) {

	//Entities
	router.GET("/:col/", H(ListEntities))
	// router.PUT("/:col/", ReplaceEntities)
	// router.DELETE("/:col/", DeleteEntities)

	// Entity
	router.GET("/:col/:id", H(RetrieveEntity))
	router.PUT("/:col/:id", H(AddEntity))
	router.DELETE("/:col/:id", H(DeleteEntity))

	// Fields
	router.GET("/:col/:id/*fieldpath", H(RetrieveField))
	router.PUT("/:col/:id/*fieldpath", H(UpdateField))
	router.DELETE("/:col/:id/*fieldpath", H(DeleteField))

}

func ListEntities(ctx *context, w http.ResponseWriter, req *http.Request) (interface{}, error) {

	// Esto debería  ser incremental ...
	col := ctx.params[0].Value
	entities, err := store.FindAll(ctx, col)
	if err != nil {
		ctx.Infof("error finding all: %v", err)
		return nil, err
	}
	return entities, nil
}

func RetrieveEntity(ctx *context, w http.ResponseWriter, req *http.Request) (interface{}, error) {

	col := ctx.params[0].Value
	id := ctx.params[1].Value
	ctx.Debugf("col: %q id: %q", col, id)
	ent, err := store.FindByID(ctx, col, id)
	if err != nil {
		return nil, err
	}
	return ent, nil
}

func AddEntity(ctx *context, w http.ResponseWriter, req *http.Request) (interface{}, error) {
	var entity map[string]interface{}
	entity, isObject := ctx.input.(map[string]interface{})
	if !isObject {
		return nil, ErrObjectExpected
	}
	col := ctx.params[0].Value
	id := ctx.params[1].Value
	ctx.Debugf("col: %q id: %q", col, id)
	entity["_id"] = id
	err := store.Save(ctx, col, entity)
	if err != nil {
		ctx.Infof("error saving entity: %v", err)
		return nil, err
	}
	w.WriteHeader(http.StatusCreated) // ? Created or no content
	return nil, nil
}

func DeleteEntity(ctx *context, w http.ResponseWriter, req *http.Request) (interface{}, error) {
	col := ctx.params[0].Value
	id := ctx.params[1].Value
	ctx.Debugf("col: %q id: %q", col, id)
	err := store.Delete(ctx, col, id)
	if err != nil {
		ctx.Infof("error deleting entity: %v", err)
		return nil, err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil, nil
}

func RetrieveField(ctx *context, w http.ResponseWriter, req *http.Request) (interface{}, error) {
	col := ctx.params[0].Value
	id := ctx.params[1].Value
	field := ctx.params[2].Value
	ctx.Debugf("col: %q id: %q, field: %q", col, id, field)
	field = cookField(field)

	value, err := store.FindField(ctx, col, id, field)
	if err != nil {
		ctx.Debugf("error finding field: %v", err)
		return nil, err
	}
	return value, nil
}

func DeleteField(ctx *context, w http.ResponseWriter, req *http.Request) (interface{}, error) {
	col := ctx.params[0].Value
	id := ctx.params[1].Value
	field := ctx.params[2].Value
	ctx.Debugf("col: %q id: %q, field: %q", col, id, field)
	field = cookField(field)

	err := store.DeleteField(ctx, col, id, field)
	if err != nil {
		ctx.Infof("error deleting field: %v", err)
		return nil, err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil, nil
}

func UpdateField(ctx *context, w http.ResponseWriter, req *http.Request) (interface{}, error) {
	col := ctx.params[0].Value
	id := ctx.params[1].Value
	field := ctx.params[2].Value
	ctx.Debugf("col: %q id: %q, field: %q", col, id, field)
	field = cookField(field)

	err := store.UpdateField(ctx, col, id, field, ctx.input)
	if err != nil {
		ctx.Infof("error updating field: %v")
		return nil, err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil, nil
}

func cookField(rawField string) string {
	return strings.Replace(strings.Trim(rawField, "/"), "/", ".", -1)
}
