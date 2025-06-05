package enums

type CreditType int

const (
	TypeCreditIn CreditType = iota + 1
	TypeCreditOut
)

func (dt CreditType) String() string {
	return [...]string{"CreditIn", "CreditOut"}[dt-1]
}

var CreditTypeLocaleKeyToId = map[string]CreditType{
	"CreditIn":  TypeCreditIn,
	"CreditOut": TypeCreditOut,
}
