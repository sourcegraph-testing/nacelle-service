package bussard

type (
	TestSimpleProcess struct {
		Value *IntWrapper `service:"value"`
	}

	TestUnsettableService struct {
		value *IntWrapper `service:"value"`
	}

	TestNonPointerField struct {
		Value IntWrapper `service:"value"`
	}

	TestOptionalServiceProcess struct {
		Value *IntWrapper `service:"value" optional:"true"`
	}

	TestBadOptionalServiceProcess struct {
		Value *IntWrapper `service:"value" optional:"yup"`
	}
)
