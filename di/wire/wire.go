package wire

import (
	"html2pdf/internal"

	"github.com/google/wire"
)

func InitApp() *internal.App {
	wire.Build(ProviderSet)
	return nil
}
