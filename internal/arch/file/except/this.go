package except

func This(value string) *Expression {
	return &Expression{
		value: value,
	}
}
