package bussard

import "fmt"

type (
	SimplePostInjectProcess struct {
		IValue *IntWrapper `service:"value"`
		FValue *FloatWrapper
	}

	ErrorPostInjectProcess struct{}

	RootInjectProcess struct {
		Services ServiceContainer         `service:"container"`
		Child    *SimplePostInjectProcess `service:"process"`
	}
)

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
