package underlying

import "go/types"

// Convert converts specified type to a type which contains only basic types.
// e.g. (MyInt is defined by "type MyInt int")
//  []MyInt -> []int
//  <-chan MyInt -> <-chan int
//  map[MyInt]MyInt -> map[int]int
//
// if ptrElem is true, Convert converts a pointer type to its elem type.
// e.g. *string -> string
func Convert(t types.Type, ptrElem bool) types.Type {
	c := &Converter{PtrElem: ptrElem}
	return c.Convert(t)
}

// Converter converts specified type to a type which contains only basic types.
type Converter struct {
	PtrElem bool
}

// Convert converts specified type to a type which contains only basic types.
// e.g. (MyInt is defined by "type MyInt int")
//  []MyInt -> []int
//  <-chan MyInt -> <-chan int
//  map[MyInt]MyInt -> map[int]int
//
// if ptrElem is true, Convert converts a pointer type to its elem type.
// e.g. *string -> string
func (c *Converter) Convert(t types.Type) types.Type {
	if t == nil {
		return nil
	}

	// basic types
	for _, bt := range types.Typ {
		if types.Identical(t, bt) {
			return t
		}
	}

	switch t := t.(type) {
	case *types.Struct:
		return c.convertStruct(t)
	case *types.Signature:
		return c.convertSignature(t)
	case *types.Pointer:
		if c.PtrElem {
			return c.Convert(t.Elem())
		}
		return types.NewPointer(c.Convert(t.Elem()))
	case *types.Array:
		return types.NewArray(c.Convert(t.Elem()), t.Len())
	case *types.Slice:
		return types.NewSlice(c.Convert(t.Elem()))
	case *types.Map:
		return types.NewMap(c.Convert(t.Key()), c.Convert(t.Elem()))
	case *types.Chan:
		return types.NewChan(t.Dir(), c.Convert(t.Elem()))
	}

	// named type
	return c.Convert(t.Underlying())
}

func (c *Converter) convertStruct(strct *types.Struct) *types.Struct {
	fields := make([]*types.Var, strct.NumFields())
	tags := make([]string, strct.NumFields())
	for i := 0; i < strct.NumFields(); i++ {
		f := strct.Field(i)
		ft := c.Convert(f.Type())
		fields[i] = types.NewField(f.Pos(), f.Pkg(), f.Name(), ft, f.Embedded())
		tags[i] = strct.Tag(i)
	}
	return types.NewStruct(fields, tags)
}

func (c *Converter) convertSignature(s *types.Signature) *types.Signature {
	params := c.convertParams(s.Params())
	results := c.convertParams(s.Results())
	return types.NewSignature(s.Recv(), params, results, s.Variadic())
}

func (c *Converter) convertParams(t *types.Tuple) *types.Tuple {
	params := make([]*types.Var, t.Len())
	for i := 0; i < t.Len(); i++ {
		p := t.At(i)
		params[i] = types.NewParam(p.Pos(), p.Pkg(), p.Name(), c.Convert(p.Type()))
	}
	return types.NewTuple(params...)
}
