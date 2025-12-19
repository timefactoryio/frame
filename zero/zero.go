package zero

type Zero interface {
	Fx
	Forge
	Element
}

type zero struct {
	Fx
	Forge
	Element
}

func NewZero(pathlessUrl, apiUrl string) Zero {
	z := &zero{
		Fx:      NewFx(pathlessUrl, apiUrl).(*fx),
		Forge:   NewForge().(*forge),
		Element: NewElement().(*element),
	}
	return z
}
