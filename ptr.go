package luar

import (
	"reflect"

	"github.com/bitwalker/gopher-lua"
)

func checkPtr(L *lua.LState, idx int) (ref reflect.Value, mt *Metatable) {
	ud := L.CheckUserData(idx)
	ref = reflect.ValueOf(ud.Value)
	if ref.Kind() != reflect.Ptr {
		L.ArgError(idx, "expecting pointer")
	}
	mt = &Metatable{LTable: ud.Metatable.(*lua.LTable)}
	return
}

func ptrIndex(L *lua.LState) int {
	_, mt := checkPtr(L, 1)
	key := L.CheckString(2)

	if fn := mt.ptrMethod(key); fn != nil {
		L.Push(fn)
		return 1
	}

	if fn := mt.method(key); fn != nil {
		L.Push(fn)
		return 1
	}

	return 0
}

func ptrPow(L *lua.LState) int {
	ref, _ := checkPtr(L, 1)
	val := L.CheckAny(2)

	if ref.IsNil() {
		L.RaiseError("cannot dereference nil pointer")
	}
	elem := ref.Elem()
	if !elem.CanSet() {
		L.RaiseError("unable to set pointer value")
	}
	value := lValueToReflect(val, elem.Type())
	elem.Set(value)
	return 1
}

func ptrUnm(L *lua.LState) int {
	ref, _ := checkPtr(L, 1)
	elem := ref.Elem()
	if !elem.CanInterface() {
		L.RaiseError("cannot interface pointer type " + elem.String())
	}
	L.Push(New(L, elem.Interface()))
	return 1
}
