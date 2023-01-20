package ir

const sample = `<h1>Hello {{customer.name}}</h1>
{% if customer.value > 100 %}
	Your account is ranked {{ customer.tier }}
{% elsif customer.value == 99 %}
	You're nearly ready to be ranked!
{% else %}
	Increase your value to get a ranking <a href="{{ customer.premium_upsell_url | tidy_url }}">here</a>`

var rep = []Node{
	String("<h1>Hello "),
	Variable("customer.name"),
	String("</h1>\n"),
	Condition([]Conditional{
		If(
			Variable("customer.value"), GT, Literal(10),
			Block([]Node{
				String("Your account is ranked "),
				Variable("customer.tier"),
				String("\n"),
			}),
		),
		If(
			Variable("customer.value"), Equals, Literal(99),
			Block([]Node{
				String("You're nearly ready to be ranked!\n"),
			}),
		),
		Else(
			String(`Increase your value to get a ranking <a href="`),
			Variable("customer.premium_upsell_url", TidyURL),
			String("\">here</a>\n"),
		),
	}),
}
