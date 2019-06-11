package service

import (
	"fmt"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type ServiceContainerSuite struct{}

func (s *ServiceContainerSuite) TestGetAndSet(t sweet.T) {
	container := NewServiceContainer()
	container.Set("a", &IntWrapper{10})
	container.Set("b", &FloatWrapper{3.14})
	container.Set("c", &IntWrapper{25})

	value1, err1 := container.Get("a")
	Expect(err1).To(BeNil())
	Expect(value1).To(Equal(&IntWrapper{10}))

	value2, err2 := container.Get("b")
	Expect(err2).To(BeNil())
	Expect(value2).To(Equal(&FloatWrapper{3.14}))

	value3, err3 := container.Get("c")
	Expect(err3).To(BeNil())
	Expect(value3).To(Equal(&IntWrapper{25}))
}

func (s *ServiceContainerSuite) TestInject(t sweet.T) {
	container := NewServiceContainer()
	container.Set("value", &IntWrapper{42})
	obj := &TestSimpleProcess{}
	err := container.Inject(obj)
	Expect(err).To(BeNil())
	Expect(obj.Value.val).To(Equal(42))
}

func (s *ServiceContainerSuite) TestInjectNonPointer(t sweet.T) {
	container := NewServiceContainer()
	container.Set("value", IntWrapper{42})
	obj := &TestSimpleNonPointer{}
	err := container.Inject(obj)
	Expect(err).To(BeNil())
	Expect(obj.Value.val).To(Equal(42))
}

func (s *ServiceContainerSuite) TestInjectAnonymous(t sweet.T) {
	container := NewServiceContainer()
	container.Set("value", &IntWrapper{42})
	obj := &TestAnonymousSimpleProcess{&TestSimpleProcess{}}
	err := container.Inject(obj)
	Expect(err).To(BeNil())
	Expect(obj.Value.val).To(Equal(42))
}

func (s *ServiceContainerSuite) TestInjectAnonymousZeroValue(t sweet.T) {
	container := NewServiceContainer()
	container.Set("value", &IntWrapper{42})
	obj := &TestAnonymousSimpleProcess{}
	err := container.Inject(obj)
	Expect(err).To(BeNil())
	Expect(obj.Value.val).To(Equal(42))
}

func (s *ServiceContainerSuite) TestInjectAnonymousNonPointer(t sweet.T) {
	container := NewServiceContainer()
	container.Set("value", &IntWrapper{42})
	obj := &TestAnonymousNonPointerSimpleProcess{}
	err := container.Inject(obj)
	Expect(err).To(BeNil())
	Expect(obj.Value.val).To(Equal(42))
}

func (s *ServiceContainerSuite) TestInjectAnonymousDeepNonPointer(t sweet.T) {
	container := NewServiceContainer()
	container.Set("value", &IntWrapper{42})
	obj := &TestAnonymousDeepNonPointerSimpleProcess{}
	err := container.Inject(obj)
	Expect(err).To(BeNil())
	Expect(obj.Value.val).To(Equal(42))
}

func (s *ServiceContainerSuite) TestInjectAnonymousZeroValueNoServiceTags(t sweet.T) {
	container := NewServiceContainer()
	container.Set("value", &IntWrapper{42})
	obj := &TestAnonymousNoServiceTags{}
	err := container.Inject(obj)
	Expect(err).To(BeNil())
	Expect(obj.IntWrapper).To(BeNil())
}

func (s *ServiceContainerSuite) TestInjectAnonymousUnexported(t sweet.T) {
	container := NewServiceContainer()
	container.Set("value", &IntWrapper{42})
	obj := &TestAnonymousUnexportedProcess{&privateProcess{}}
	err := container.Inject(obj)
	Expect(err).To(BeNil())
	Expect(obj.privateProcess.Value).To(BeNil())
}

func (s *ServiceContainerSuite) TestInjectNonStruct(t sweet.T) {
	container := NewServiceContainer()
	obj := func() error { return nil }
	err := container.Inject(obj)
	Expect(err).To(BeNil())
}

func (s *ServiceContainerSuite) TestInjectMissingService(t sweet.T) {
	container := NewServiceContainer()
	obj := &TestSimpleProcess{}
	err := container.Inject(obj)
	Expect(err).To(MatchError("no service registered to key `value`"))
}

func (s *ServiceContainerSuite) TestInjectBadType(t sweet.T) {
	container := NewServiceContainer()
	container.Set("value", &FloatWrapper{3.14})
	obj := &TestSimpleProcess{}
	err := container.Inject(obj)
	Expect(err).To(MatchError("field 'Value' cannot be assigned a value of type *service.FloatWrapper"))
}

func (s *ServiceContainerSuite) TestInjectNil(t sweet.T) {
	container := NewServiceContainer()
	container.Set("value", nil)
	obj := &TestNonPointerField{}
	err := container.Inject(obj)
	Expect(err).To(MatchError("field 'Value' cannot be assigned a value of type nil"))
}

func (s *ServiceContainerSuite) TestInjectOptional(t sweet.T) {
	container := NewServiceContainer()
	obj := &TestOptionalServiceProcess{}
	err := container.Inject(obj)
	Expect(err).To(BeNil())
	Expect(obj.Value).To(BeNil())

	container.Set("value", &IntWrapper{42})
	err = container.Inject(obj)
	Expect(err).To(BeNil())
	Expect(obj.Value.val).To(Equal(42))
}

func (s *ServiceContainerSuite) TestInjectBadOptional(t sweet.T) {
	container := NewServiceContainer()
	obj := &TestBadOptionalServiceProcess{}
	err := container.Inject(obj)
	Expect(err).To(MatchError("field 'Value' has an invalid optional tag"))
}

func (s *ServiceContainerSuite) TestUnsettableFields(t sweet.T) {
	container := NewServiceContainer()
	container.Set("value", &IntWrapper{42})
	err := container.Inject(&TestUnsettableService{})
	Expect(err).To(MatchError("field 'value' can not be set - it may be unexported"))
}

func (s *ServiceContainerSuite) TestPostInject(t sweet.T) {
	container := NewServiceContainer()
	obj := &SimplePostInjectProcess{}
	container.Set("value", &IntWrapper{42})
	err := container.Inject(obj)
	Expect(err).To(BeNil())
	Expect(obj.FValue.val).To(Equal(42.0))
}

func (s *ServiceContainerSuite) TestPostInjectError(t sweet.T) {
	container := NewServiceContainer()
	obj := &ErrorPostInjectProcess{}
	container.Set("value", &IntWrapper{42})
	err := container.Inject(obj)
	Expect(err).To(MatchError("utoh"))
}

func (s *ServiceContainerSuite) TestPostInjectChain(t sweet.T) {
	container := NewServiceContainer()
	obj := &RootInjectProcess{}
	process := &SimplePostInjectProcess{}

	container.Set("value", &IntWrapper{42})
	container.Set("process", process)
	container.Set("services", container)

	err := container.Inject(obj)
	Expect(err).To(BeNil())
	Expect(process.FValue.val).To(Equal(42.0))
}

func (s *ServiceContainerSuite) TestDuplicateRegistration(t sweet.T) {
	container := NewServiceContainer()
	err1 := container.Set("dup", struct{}{})
	err2 := container.Set("dup", struct{}{})
	Expect(err1).To(BeNil())
	Expect(err2).To(MatchError("duplicate service key `dup`"))
}

func (s *ServiceContainerSuite) TestGetUnregisteredKey(t sweet.T) {
	container := NewServiceContainer()
	_, err := container.Get("unregistered")
	Expect(err).To(MatchError("no service registered to key `unregistered`"))
}

func (s *ServiceContainerSuite) TestMustSetPanics(t sweet.T) {
	Expect(func() {
		container := NewServiceContainer()
		container.MustSet("unregistered", struct{}{})
		container.MustSet("unregistered", struct{}{})
	}).To(Panic())
}

func (s *ServiceContainerSuite) TestMustGetPanics(t sweet.T) {
	Expect(func() {
		NewServiceContainer().MustGet("unregistered")
	}).To(Panic())
}

//
// Fixtures

type (
	IntWrapper struct {
		val int
	}

	FloatWrapper struct {
		val float64
	}

	TestSimpleNonPointer struct {
		Value IntWrapper `service:"value"`
	}

	TestSimpleProcess struct {
		Value *IntWrapper `service:"value"`
	}

	TestAnonymousSimpleProcess struct {
		*TestSimpleProcess
	}

	TestAnonymousNonPointerSimpleProcess struct {
		TestSimpleProcess
	}

	TestAnonymousDeepNonPointerSimpleProcess struct {
		TestSimpleProcess
	}

	TestAnonymousNoServiceTags struct {
		*IntWrapper
	}

	TestAnonymousUnexportedProcess struct {
		*privateProcess
	}

	privateProcess struct {
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

	SimplePostInjectProcess struct {
		IValue *IntWrapper `service:"value"`
		FValue *FloatWrapper
	}

	ErrorPostInjectProcess struct{}

	RootInjectProcess struct {
		Services ServiceContainer         `service:"services"`
		Child    *SimplePostInjectProcess `service:"process"`
	}
)

//
// Fixture Methods

func (p *SimplePostInjectProcess) PostInject() error {
	p.FValue = &FloatWrapper{float64(p.IValue.val)}
	return nil
}

func (p *ErrorPostInjectProcess) PostInject() error {
	return fmt.Errorf("utoh")
}

func (p *RootInjectProcess) PostInject() error {
	return p.Services.Inject(p.Child)
}
