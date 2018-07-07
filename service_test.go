package bussard

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type ServiceSuite struct{}

func (s *ServiceSuite) TestGetAndSet(t sweet.T) {
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

func (s *ServiceSuite) TestInject(t sweet.T) {
	container := NewServiceContainer()
	container.Set("value", &IntWrapper{42})
	obj := &TestSimpleProcess{}
	err := container.Inject(obj)
	Expect(err).To(BeNil())
	Expect(obj.Value.val).To(Equal(42))
}

func (s *ServiceSuite) TestInjectNonStruct(t sweet.T) {
	container := NewServiceContainer()
	obj := func() error { return nil }
	err := container.Inject(obj)
	Expect(err).To(BeNil())
}

func (s *ServiceSuite) TestInjectMissingService(t sweet.T) {
	container := NewServiceContainer()
	obj := &TestSimpleProcess{}
	err := container.Inject(obj)
	Expect(err).To(MatchError("no service registered to key `value`"))
}

func (s *ServiceSuite) TestInjectBadType(t sweet.T) {
	container := NewServiceContainer()
	container.Set("value", &FloatWrapper{3.14})
	obj := &TestSimpleProcess{}
	err := container.Inject(obj)
	Expect(err).To(MatchError("field 'Value' cannot be assigned a value of type *bussard.FloatWrapper"))
}

func (s *ServiceSuite) TestInjectNil(t sweet.T) {
	container := NewServiceContainer()
	container.Set("value", nil)
	obj := &TestNonPointerField{}
	err := container.Inject(obj)
	Expect(err).To(MatchError("field 'Value' cannot be assigned a value of type nil"))
}

func (s *ServiceSuite) TestInjectOptional(t sweet.T) {
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

func (s *ServiceSuite) TestInjectBadOptional(t sweet.T) {
	container := NewServiceContainer()
	obj := &TestBadOptionalServiceProcess{}
	err := container.Inject(obj)
	Expect(err).To(MatchError("field 'Value' has an invalid optional tag"))
}

func (s *ServiceSuite) TestUnsettableFields(t sweet.T) {
	container := NewServiceContainer()
	container.Set("value", &IntWrapper{42})
	err := container.Inject(&TestUnsettableService{})
	Expect(err).To(MatchError("field 'value' can not be set - it may be unexported"))
}

func (s *ServiceSuite) TestPostInject(t sweet.T) {
	container := NewServiceContainer()
	obj := &SimplePostInjectProcess{}
	container.Set("value", &IntWrapper{42})
	err := container.Inject(obj)
	Expect(err).To(BeNil())
	Expect(obj.FValue.val).To(Equal(42.0))
}

func (s *ServiceSuite) TestPostInjectError(t sweet.T) {
	container := NewServiceContainer()
	obj := &ErrorPostInjectProcess{}
	container.Set("value", &IntWrapper{42})
	err := container.Inject(obj)
	Expect(err).To(MatchError("utoh"))
}

func (s *ServiceSuite) TestPostInjectChain(t sweet.T) {
	container := NewServiceContainer()
	obj := &RootInjectProcess{}
	process := &SimplePostInjectProcess{}

	container.Set("value", &IntWrapper{42})
	container.Set("process", process)
	container.Set("container", container)

	err := container.Inject(obj)
	Expect(err).To(BeNil())
	Expect(process.FValue.val).To(Equal(42.0))
}

func (s *ServiceSuite) TestDuplicateRegistration(t sweet.T) {
	container := NewServiceContainer()
	err1 := container.Set("dup", struct{}{})
	err2 := container.Set("dup", struct{}{})
	Expect(err1).To(BeNil())
	Expect(err2).To(MatchError("duplicate service key `dup`"))
}

func (s *ServiceSuite) TestGetUnregisteredKey(t sweet.T) {
	container := NewServiceContainer()
	_, err := container.Get("unregistered")
	Expect(err).To(MatchError("no service registered to key `unregistered`"))
}

func (s *ServiceSuite) TestMustSetPanics(t sweet.T) {
	Expect(func() {
		container := NewServiceContainer()
		container.MustSet("unregistered", struct{}{})
		container.MustSet("unregistered", struct{}{})
	}).To(Panic())
}

func (s *ServiceSuite) TestMustGetPanics(t sweet.T) {
	Expect(func() {
		NewServiceContainer().MustGet("unregistered")
	}).To(Panic())
}
