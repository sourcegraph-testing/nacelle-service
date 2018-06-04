package bussard

type IntWrapper struct {
	val int
}

type FloatWrapper struct {
	val float64
}

type TestSimpleProcess struct {
	Value *IntWrapper `service:"value"`
}

type TestUnsettableService struct {
	value *IntWrapper `service:"value"`
}

type TestNonPointerField struct {
	Value IntWrapper `service:"value"`
}

type TestOptionalServiceProcess struct {
	Value *IntWrapper `service:"value" optional:"true"`
}

type TestBadOptionalServiceProcess struct {
	Value *IntWrapper `service:"value" optional:"yup"`
}
