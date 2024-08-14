package templates

import (
	"fmt"
)

type ExchangeRateTemplate struct {
	CreatedAt    string
	ExchangeRate string
}

func (ExchangeRateTemplate) Template() string {
	return fmt.Sprintf("%s %s %s", StartTempStyle,
		`<h3>Hello!</h3>
	<p><b>Date:</b> {{.CreatedAt}}.</p>
	<p>The current exchange rate for USD to UAH is: {{.ExchangeRate}} ₴.</p>
	<br/>

	<p style="font-size: 14px; color: #696969;">
	Kind regards, <br>
	Exchange Rates Team</p>`, EndTempStyle)
}
