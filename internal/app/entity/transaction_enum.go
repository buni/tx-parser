// Code generated by "enumer -type=TokenType -trimprefix=TokenType -transform=snake -output=transaction_enum.go -json -sql -text"; DO NOT EDIT.

package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

const _TokenTypeName = "invliadeth"

var _TokenTypeIndex = [...]uint8{0, 7, 10}

const _TokenTypeLowerName = "invliadeth"

func (i TokenType) String() string {
	if i < 0 || i >= TokenType(len(_TokenTypeIndex)-1) {
		return fmt.Sprintf("TokenType(%d)", i)
	}
	return _TokenTypeName[_TokenTypeIndex[i]:_TokenTypeIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _TokenTypeNoOp() {
	var x [1]struct{}
	_ = x[TokenTypeInvliad-(0)]
	_ = x[TokenTypeETH-(1)]
}

var _TokenTypeValues = []TokenType{TokenTypeInvliad, TokenTypeETH}

var _TokenTypeNameToValueMap = map[string]TokenType{
	_TokenTypeName[0:7]:       TokenTypeInvliad,
	_TokenTypeLowerName[0:7]:  TokenTypeInvliad,
	_TokenTypeName[7:10]:      TokenTypeETH,
	_TokenTypeLowerName[7:10]: TokenTypeETH,
}

var _TokenTypeNames = []string{
	_TokenTypeName[0:7],
	_TokenTypeName[7:10],
}

// TokenTypeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func TokenTypeString(s string) (TokenType, error) {
	if val, ok := _TokenTypeNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _TokenTypeNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to TokenType values", s)
}

// TokenTypeValues returns all values of the enum
func TokenTypeValues() []TokenType {
	return _TokenTypeValues
}

// TokenTypeStrings returns a slice of all String values of the enum
func TokenTypeStrings() []string {
	strs := make([]string, len(_TokenTypeNames))
	copy(strs, _TokenTypeNames)
	return strs
}

// IsATokenType returns "true" if the value is listed in the enum definition. "false" otherwise
func (i TokenType) IsATokenType() bool {
	for _, v := range _TokenTypeValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for TokenType
func (i TokenType) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for TokenType
func (i *TokenType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("TokenType should be a string, got %s", data)
	}

	var err error
	*i, err = TokenTypeString(s)
	return err
}

// MarshalText implements the encoding.TextMarshaler interface for TokenType
func (i TokenType) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for TokenType
func (i *TokenType) UnmarshalText(text []byte) error {
	var err error
	*i, err = TokenTypeString(string(text))
	return err
}

func (i TokenType) Value() (driver.Value, error) {
	return i.String(), nil
}

func (i *TokenType) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	case fmt.Stringer:
		str = v.String()
	default:
		return fmt.Errorf("invalid value of TokenType: %[1]T(%[1]v)", value)
	}

	val, err := TokenTypeString(str)
	if err != nil {
		return err
	}

	*i = val
	return nil
}
