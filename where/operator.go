package where

type operator uint8

const (
	opEqual operator = iota
	opGreater
	opLess
	opLike
	opIn
	opBetween
	opIsNull
)

var operators [14]string = [14]string{
	"`%s` = %s",
	"`%s` > %s",
	"`%s` < %s",
	"`%s` like %s",
	"`%s` in (%s)",
	"`%s` between %s and %s",
	"`%s` is null",
	"`%s` <> %s",
	"`%s` <= %s",
	"`%s` >= %s",
	"`%s` not like %s",
	"`%s` not in (%s)",
	"`%s` not between %s and %s",
	"`%s` is not null",
}

var operatorType [7]int = [7]int{
	dBool & dDate & dFloat & dDouble & dInt & dLong & dText,
	dDate & dFloat & dDouble & dInt & dLong & dText,
	dDate & dFloat & dDouble & dInt & dLong & dText,
	dText,
	dDate & dFloat & dDouble & dInt & dLong & dText,
	dDate & dFloat & dDouble & dInt & dLong,
	dBool & dDate & dFloat & dDouble & dInt & dLong & dText,
}
