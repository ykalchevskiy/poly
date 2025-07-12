package poly

// TypeList is a recursive type that represents a list of types.
// It is used to provide a variadic number of types to a generic function.
type TypeList[First TypeName, Rest Types] struct{}

func (TypeList[First, Rest]) Types() []Type {
	var r Rest

	return append([]Type{NewType[First]()}, r.Types()...)
}

// TypeListLast is the terminating type for TypeList, representing the end of the list.
type TypeListLast struct{}

func (TypeListLast) Types() []Type {
	return nil
}

// Types1 represents a list containing a single type.
type Types1[T1 TypeName] struct{}

func (Types1[T1]) Types() []Type {
	return []Type{
		NewType[T1](),
	}
}

// Types2 represents a list containing two types.
type Types2[T1, T2 TypeName] struct{}

func (Types2[T1, T2]) Types() []Type {
	return []Type{
		NewType[T1](),
		NewType[T2](),
	}
}

// Types3 represents a list containing three types.
type Types3[T1, T2, T3 TypeName] struct{}

func (Types3[T1, T2, T3]) Types() []Type {
	return []Type{
		NewType[T1](),
		NewType[T2](),
		NewType[T3](),
	}
}

// Types4 represents a list containing four types.
type Types4[T1, T2, T3, T4 TypeName] struct{}

func (Types4[T1, T2, T3, T4]) Types() []Type {
	return []Type{
		NewType[T1](),
		NewType[T2](),
		NewType[T3](),
		NewType[T4](),
	}
}

// Types5 represents a list containing five types.
type Types5[T1, T2, T3, T4, T5 TypeName] struct{}

func (Types5[T1, T2, T3, T4, T5]) Types() []Type {
	return []Type{
		NewType[T1](),
		NewType[T2](),
		NewType[T3](),
		NewType[T4](),
		NewType[T5](),
	}
}

// Types6 represents a list containing six types.
type Types6[T1, T2, T3, T4, T5, T6 TypeName] struct{}

func (Types6[T1, T2, T3, T4, T5, T6]) Types() []Type {
	return []Type{
		NewType[T1](),
		NewType[T2](),
		NewType[T3](),
		NewType[T4](),
		NewType[T5](),
		NewType[T6](),
	}
}

// Types7 represents a list containing seven types.
type Types7[T1, T2, T3, T4, T5, T6, T7 TypeName] struct{}

func (Types7[T1, T2, T3, T4, T5, T6, T7]) Types() []Type {
	return []Type{
		NewType[T1](),
		NewType[T2](),
		NewType[T3](),
		NewType[T4](),
		NewType[T5](),
		NewType[T6](),
		NewType[T7](),
	}
}

// Types8 represents a list containing eight types.
type Types8[T1, T2, T3, T4, T5, T6, T7, T8 TypeName] struct{}

func (Types8[T1, T2, T3, T4, T5, T6, T7, T8]) Types() []Type {
	return []Type{
		NewType[T1](),
		NewType[T2](),
		NewType[T3](),
		NewType[T4](),
		NewType[T5](),
		NewType[T6](),
		NewType[T7](),
		NewType[T8](),
	}
}

// Types9 represents a list containing nine types.
type Types9[T1, T2, T3, T4, T5, T6, T7, T8, T9 TypeName] struct{}

func (Types9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) Types() []Type {
	return []Type{
		NewType[T1](),
		NewType[T2](),
		NewType[T3](),
		NewType[T4](),
		NewType[T5](),
		NewType[T6](),
		NewType[T7](),
		NewType[T8](),
		NewType[T9](),
	}
}
