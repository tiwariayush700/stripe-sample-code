package model

import (
	`bytes`
	`encoding/json`
	`fmt`
	`math`
	`strconv`
	`strings`

	`github.com/Rhymond/go-money`
	`go.mongodb.org/mongo-driver/bson`
	`go.mongodb.org/mongo-driver/bson/bsontype`
	`go.mongodb.org/mongo-driver/x/bsonx/bsoncore`
)

type Money struct {
	money.Money
}

// NewMoney creates and returns new instance of Money.
func NewMoney(amount int64, code string) Money {
	return Money{Money: *money.New(amount, code)}
}

// ParseMoneyString converts string to Money, eg. "45.67", "USD" --> Money{}
func ParseMoneyString(stringValue string, curr string) Money {
	curr = strings.ToUpper(curr)
	value, _ := strconv.ParseFloat(stringValue, 64)
	fraction := math.Pow(10, float64(money.GetCurrency(curr).Fraction))
	intValue := int64(value * fraction)
	return NewMoney(intValue, curr)
}

// Split returns slice of Money structs with split Self value in given number.
func (m Money) Split(n int) ([]Money, error) {
	tmp, err := m.Money.Split(n)
	if err != nil {
		return nil, err
	}

	split := make([]Money, n, n)
	for i, m2 := range tmp {
		split[i] = NewMoney(m2.Amount(), m2.Currency().Code)
	}

	return split, nil

}

// Subtract returns new Money struct with value representing difference of Self and Other Money.
func (m Money) Subtract(om Money) (Money, error) {
	tmp, err := m.Money.Subtract(&om.Money)
	if err != nil {
		return Money{}, err
	}

	return NewMoney(tmp.Amount(), tmp.Currency().Code), nil
}

// Add returns new Money struct with value representing sum of Self and Other Money.
func (m Money) Add(om Money) (Money, error) {
	tmp, err := m.Money.Add(&om.Money)
	if err != nil {
		return Money{}, err
	}

	return NewMoney(tmp.Amount(), tmp.Currency().Code), nil
}

// Multiply returns new Money struct with value representing Self multiplied value by multiplier.
func (m Money) Multiply(mul int64) Money {
	tmp := m.Money.Multiply(mul)
	return NewMoney(tmp.Amount(), tmp.Currency().Code)
}

// MarshalBSONValue overrides BSON marshal behavior
func (m Money) MarshalBSONValue() (bsontype.Type, []byte, error) {

	idx, doc := bsoncore.AppendDocumentStart(nil)
	doc = append(doc, bsoncore.AppendStringElement(nil, "currency", m.Currency().Code)...)
	doc = append(doc, bsoncore.AppendInt64Element(nil, "amount", m.Amount())...)
	doc, _ = bsoncore.AppendDocumentEnd(doc, idx)

	return bsontype.EmbeddedDocument, doc, nil
}

// UnmarshalBSONValue overrides BSON unmarshal behavior
func (m *Money) UnmarshalBSONValue(bsonType bsontype.Type, data []byte) error {

	type tmp struct {
		Currency string `bson:"currency"`
		Amount   int64  `bson:"amount"`
	}

	var t tmp
	err := bson.Unmarshal(data, &t)
	if err != nil {
		return err
	}

	o := money.New(t.Amount, t.Currency)
	*m = Money{*o}

	return nil
}

// MarshalJSON overrides JSON marshal behavior
func (m Money) MarshalJSON() ([]byte, error) {
	f := "%." + strconv.Itoa(m.Currency().Fraction) + "f"
	format := fmt.Sprintf(`{"amount": %s, `, f) + `"currency": "%s"}`
	buff := bytes.NewBufferString(fmt.Sprintf(format, m.AsMajorUnits(), m.Currency().Code))
	return buff.Bytes(), nil
}

// UnmarshalJSON is implementation of json.Unmarshaller
func (m *Money) UnmarshalJSON(data []byte) error {

	type tmp struct {
		Currency string  `bson:"currency"`
		Amount   float64 `bson:"amount"`
	}

	var t tmp
	err := json.Unmarshal(data, &t)
	if err != nil {
		return err
	}

	m1 := ParseMoneyString(fmt.Sprintf("%f", t.Amount), t.Currency)

	*m = m1

	return nil
}
