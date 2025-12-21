package zero

import (
	"io"
	"net/http"
	"os"
	"strings"
)

type Zero struct {
	Forge
	Element
}

func NewZero() *Zero {
	z := &Zero{
		Forge:   NewForge().(*forge),
		Element: NewElement().(*element),
	}
	return z
}

func (z *Zero) BuildFromFile(html, class string) *One {
	file := z.HTML(z.ToString(html))
	final := z.Builder(class, file)
	return final
}

func (z *Zero) ToBytes(input string) []byte {
	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		resp, err := http.Get(input)
		if err != nil {
			return nil
		}
		defer resp.Body.Close()
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil
		}
		return b
	}
	b, err := os.ReadFile(input)
	if err != nil {
		return nil
	}
	return b
}

func (z *Zero) ToString(input string) string {
	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		resp, err := http.Get(input)
		if err != nil {
			return ""
		}
		defer resp.Body.Close()
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return ""
		}
		return string(b)
	}
	b, err := os.ReadFile(input)
	if err != nil {
		return ""
	}
	return string(b)
}
