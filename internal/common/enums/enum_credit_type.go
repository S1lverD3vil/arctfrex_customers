package enums

type CreditType int

const (
	TypeCreditDefault CreditType = iota
	TypeCreditIn
	TypeCreditOut
)

func (dt CreditType) String() string {
	return [...]string{"CreditIn", "CreditOut"}[dt-1]
}

var CreditTypeLocaleKeyToId = map[string]CreditType{
	"CreditIn":  TypeCreditIn,
	"CreditOut": TypeCreditOut,
}
