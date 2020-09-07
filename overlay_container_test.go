package service

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type OverlayContainerSuite struct{}

func (s *OverlayContainerSuite) TestGet(t sweet.T) {
	container := NewServiceContainer()
	container.Set("a", &IntWrapper{10})
	container.Set("b", &IntWrapper{20})
	container.Set("c", &IntWrapper{30})

	overlay := Overlay(container, map[string]interface{}{
		"a": &IntWrapper{40},
		"d": &IntWrapper{50},
	})

	value1, err1 := overlay.Get("a")
	Expect(err1).To(BeNil())
	Expect(value1).To(Equal(&IntWrapper{40}))

	value2, err2 := overlay.Get("b")
	Expect(err2).To(BeNil())
	Expect(value2).To(Equal(&IntWrapper{20}))

	value3, err3 := overlay.Get("c")
	Expect(err3).To(BeNil())
	Expect(value3).To(Equal(&IntWrapper{30}))

	value4, err4 := overlay.Get("d")
	Expect(err4).To(BeNil())
	Expect(value4).To(Equal(&IntWrapper{50}))
}

func (s *OverlayContainerSuite) TestInject(t sweet.T) {
	container := NewServiceContainer()
	container.Set("a", &IntWrapper{10})
	container.Set("b", &IntWrapper{20})
	container.Set("c", &IntWrapper{30})

	overlay := Overlay(container, map[string]interface{}{
		"a": &IntWrapper{40},
		"d": &IntWrapper{50},
	})

	obj := &TestOverlayProcess{}
	err := overlay.Inject(obj)
	Expect(err).To(BeNil())
	Expect(obj.A.val).To(Equal(40))
	Expect(obj.B.val).To(Equal(20))
	Expect(obj.C.val).To(Equal(30))
	Expect(obj.D.val).To(Equal(50))
}

//
// Fixtures

type (
	TestOverlayProcess struct {
		A *IntWrapper `service:"a"`
		B *IntWrapper `service:"b"`
		C *IntWrapper `service:"c"`
		D *IntWrapper `service:"d"`
	}
)
