package poly

type TypeList[First TypeName, Rest Types] struct{}

func (TypeList[First, Rest]) Types() []Type {
	var r Rest
	return append([]Type{NewType[First]()}, r.Types()...)
}

type TypeListLast struct{}

func (TypeListLast) Types() []Type {
	return nil
}

type Types1[T1 TypeName] struct{}

func (Types1[T1]) Types() []Type {
	return []Type{
		NewType[T1](),
	}
}

type Types2[T1, T2 TypeName] struct{}

func (Types2[T1, T2]) Types() []Type {
	return []Type{
		NewType[T1](),
		NewType[T2](),
	}
}

type Types3[T1, T2, T3 TypeName] struct{}

func (Types3[T1, T2, T3]) Types() []Type {
	return []Type{
		NewType[T1](),
		NewType[T2](),
		NewType[T3](),
	}
}

type Types4[T1, T2, T3, T4 TypeName] struct{}

func (Types4[T1, T2, T3, T4]) Types() []Type {
	return []Type{
		NewType[T1](),
		NewType[T2](),
		NewType[T3](),
		NewType[T4](),
	}
}


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
