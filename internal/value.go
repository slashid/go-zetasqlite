package internal

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/base64"
	"fmt"
	"io"
	"math/big"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-json"
	ast "github.com/goccy/go-zetasql/resolved_ast"
	"github.com/goccy/go-zetasql/types"
)

type Value interface {
	Add(Value) (Value, error)
	Sub(Value) (Value, error)
	Mul(Value) (Value, error)
	Div(Value) (Value, error)
	EQ(Value) (bool, error)
	GT(Value) (bool, error)
	GTE(Value) (bool, error)
	LT(Value) (bool, error)
	LTE(Value) (bool, error)
	ToInt64() (int64, error)
	ToString() (string, error)
	ToBytes() ([]byte, error)
	ToFloat64() (float64, error)
	ToBool() (bool, error)
	ToArray() (*ArrayValue, error)
	ToStruct() (*StructValue, error)
	ToJSON() (string, error)
	ToTime() (time.Time, error)
	ToRat() (*big.Rat, error)
	Marshal() (string, error)
	Format(verb rune) string
	Interface() interface{}
}

type IntValue int64

func (iv IntValue) Add(v Value) (Value, error) {
	v2, err := v.ToInt64()
	if err != nil {
		return nil, err
	}
	return ValueOf(int64(iv) + v2)
}

func (iv IntValue) Sub(v Value) (Value, error) {
	v2, err := v.ToInt64()
	if err != nil {
		return nil, err
	}
	return ValueOf(int64(iv) - v2)
}

func (iv IntValue) Mul(v Value) (Value, error) {
	v2, err := v.ToInt64()
	if err != nil {
		return nil, err
	}
	return ValueOf(int64(iv) * v2)
}

func (iv IntValue) Div(v Value) (Value, error) {
	v2, err := v.ToInt64()
	if err != nil {
		return nil, err
	}
	if v2 == 0 {
		return nil, fmt.Errorf("zero divided error ( %d / 0 )", iv)
	}
	return ValueOf(int64(iv) / v2)
}

func (iv IntValue) EQ(v Value) (bool, error) {
	v2, err := v.ToInt64()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to int64", v)
	}
	return int64(iv) == v2, nil
}

func (iv IntValue) GT(v Value) (bool, error) {
	v2, err := v.ToInt64()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to int64", v)
	}
	return int64(iv) > v2, nil
}

func (iv IntValue) GTE(v Value) (bool, error) {
	v2, err := v.ToInt64()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to int64", v)
	}
	return int64(iv) >= v2, nil
}

func (iv IntValue) LT(v Value) (bool, error) {
	v2, err := v.ToInt64()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to int64", v)
	}
	return int64(iv) < v2, nil
}

func (iv IntValue) LTE(v Value) (bool, error) {
	v2, err := v.ToInt64()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to int64", v)
	}
	return int64(iv) <= v2, nil
}

func (iv IntValue) ToInt64() (int64, error) {
	return int64(iv), nil
}

func (iv IntValue) ToString() (string, error) {
	return fmt.Sprint(iv), nil
}

func (iv IntValue) ToBytes() ([]byte, error) {
	return []byte(fmt.Sprint(iv)), nil
}

func (iv IntValue) ToFloat64() (float64, error) {
	return float64(iv), nil
}

func (iv IntValue) ToBool() (bool, error) {
	switch iv {
	case 0:
		return false, nil
	case 1:
		return true, nil
	default:
		return false, fmt.Errorf("failed to convert %d to bool type", iv)
	}
}

func (iv IntValue) ToArray() (*ArrayValue, error) {
	return nil, fmt.Errorf("failed to convert %d to array type", iv)
}

func (iv IntValue) ToStruct() (*StructValue, error) {
	return nil, fmt.Errorf("failed to convert %d to struct type", iv)
}

func (iv IntValue) ToJSON() (string, error) {
	return fmt.Sprint(iv), nil
}

func (iv IntValue) ToTime() (time.Time, error) {
	return time.Time{}, fmt.Errorf("failed to convert %d to time.Time type", iv)
}

func (iv IntValue) ToRat() (*big.Rat, error) {
	r := new(big.Rat)
	r.SetInt64(int64(iv))
	return r, nil
}

func (iv IntValue) Marshal() (string, error) {
	return fmt.Sprint(iv), nil
}

func (iv IntValue) Format(verb rune) string {
	return fmt.Sprint(iv)
}

func (iv IntValue) Interface() interface{} {
	return int64(iv)
}

type StringValue string

func (sv StringValue) Add(v Value) (Value, error) {
	v2, err := v.ToString()
	if err != nil {
		return nil, err
	}
	return ValueOf(string(sv) + v2)
}

func (sv StringValue) Sub(v Value) (Value, error) {
	return nil, fmt.Errorf("sub operation is unsupported for string %v", sv)
}

func (sv StringValue) Mul(v Value) (Value, error) {
	return nil, fmt.Errorf("mul operation is unsupported for string %v", sv)
}

func (sv StringValue) Div(v Value) (Value, error) {
	return nil, fmt.Errorf("div operation is unsupported for string %v", sv)
}

func (sv StringValue) EQ(v Value) (bool, error) {
	v2, err := v.ToString()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to string", v)
	}
	return string(sv) == v2, nil
}

func (sv StringValue) GT(v Value) (bool, error) {
	v2, err := v.ToString()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to string", v)
	}
	return string(sv) > v2, nil
}

func (sv StringValue) GTE(v Value) (bool, error) {
	v2, err := v.ToString()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to string", v)
	}
	return string(sv) >= v2, nil
}

func (sv StringValue) LT(v Value) (bool, error) {
	v2, err := v.ToString()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to string", v)
	}
	return string(sv) < v2, nil
}

func (sv StringValue) LTE(v Value) (bool, error) {
	v2, err := v.ToString()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to string", v)
	}
	return string(sv) <= v2, nil
}

func (sv StringValue) ToInt64() (int64, error) {
	if sv == "" {
		return 0, nil
	}
	return strconv.ParseInt(string(sv), 10, 64)
}

func (sv StringValue) ToString() (string, error) {
	return string(sv), nil
}

func (sv StringValue) ToBytes() ([]byte, error) {
	return []byte(string(sv)), nil
}

func (sv StringValue) ToFloat64() (float64, error) {
	if sv == "" {
		return 0, nil
	}
	return strconv.ParseFloat(string(sv), 64)
}

func (sv StringValue) ToBool() (bool, error) {
	if sv == "" {
		return false, nil
	}
	return strconv.ParseBool(string(sv))
}

func (sv StringValue) ToArray() (*ArrayValue, error) {
	if sv == "" {
		return nil, nil
	}
	return nil, fmt.Errorf("failed to convert array from string: %v", sv)
}

func (sv StringValue) ToStruct() (*StructValue, error) {
	if sv == "" {
		return nil, nil
	}
	return nil, fmt.Errorf("failed to convert struct from string: %v", sv)
}

func (sv StringValue) ToJSON() (string, error) {
	return strconv.Quote(string(sv)), nil
}

func (sv StringValue) ToTime() (time.Time, error) {
	switch {
	case isDate(string(sv)):
		return parseDate(string(sv))
	}
	return time.Time{}, fmt.Errorf("failed to convert %s to time.Time type", sv)
}

func (sv StringValue) ToRat() (*big.Rat, error) {
	r := new(big.Rat)
	r.SetString(string(sv))
	return r, nil
}

func (sv StringValue) Marshal() (string, error) {
	return strconv.Quote(string(sv)), nil
}

func (sv StringValue) Format(verb rune) string {
	switch verb {
	case 't':
		return string(sv)
	case 'T':
		return strconv.Quote(string(sv))
	}
	return string(sv)
}

func (sv StringValue) Interface() interface{} {
	return string(sv)
}

type BytesValue []byte

func (bv BytesValue) Add(v Value) (Value, error) {
	v2, err := v.ToBytes()
	if err != nil {
		return nil, err
	}
	return BytesValue(append([]byte(bv), v2...)), nil
}

func (bv BytesValue) Sub(v Value) (Value, error) {
	return nil, fmt.Errorf("sub operation is unsupported for bytes %v", bv)
}

func (bv BytesValue) Mul(v Value) (Value, error) {
	return nil, fmt.Errorf("mul operation is unsupported for bytes %v", bv)
}

func (bv BytesValue) Div(v Value) (Value, error) {
	return nil, fmt.Errorf("div operation is unsupported for bytes %v", bv)
}

func (bv BytesValue) EQ(v Value) (bool, error) {
	v2, err := v.ToBytes()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to bytes", v)
	}
	return bytes.Equal([]byte(bv), v2), nil
}

func (bv BytesValue) GT(v Value) (bool, error) {
	v2, err := v.ToBytes()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to bytes", v)
	}
	return string(bv) > string(v2), nil
}

func (bv BytesValue) GTE(v Value) (bool, error) {
	v2, err := v.ToBytes()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to bytes", v)
	}
	return string(bv) >= string(v2), nil
}

func (bv BytesValue) LT(v Value) (bool, error) {
	v2, err := v.ToBytes()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to bytes", v)
	}
	return string(bv) < string(v2), nil
}

func (bv BytesValue) LTE(v Value) (bool, error) {
	v2, err := v.ToBytes()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to bytes", v)
	}
	return string(bv) <= string(v2), nil
}

func (bv BytesValue) ToInt64() (int64, error) {
	if len(bv) == 0 {
		return 0, nil
	}
	return strconv.ParseInt(string(bv), 10, 64)
}

func (bv BytesValue) ToString() (string, error) {
	return base64.StdEncoding.EncodeToString([]byte(bv)), nil
}

func (bv BytesValue) ToBytes() ([]byte, error) {
	return []byte(bv), nil
}

func (bv BytesValue) ToFloat64() (float64, error) {
	if len(bv) == 0 {
		return 0, nil
	}
	return strconv.ParseFloat(string(bv), 64)
}

func (bv BytesValue) ToBool() (bool, error) {
	if len(bv) == 0 {
		return false, nil
	}
	return strconv.ParseBool(string(bv))
}

func (bv BytesValue) ToArray() (*ArrayValue, error) {
	return nil, fmt.Errorf("failed to convert array from bytes: %v", bv)
}

func (bv BytesValue) ToStruct() (*StructValue, error) {
	return nil, fmt.Errorf("failed to convert struct from bytes: %v", bv)
}

func (bv BytesValue) ToJSON() (string, error) {
	v, err := bv.ToString()
	if err != nil {
		return "", err
	}
	return strconv.Quote(v), nil
}

func (bv BytesValue) ToTime() (time.Time, error) {
	v := string(bv)
	switch {
	case isDate(v):
		return parseDate(v)
	}
	return time.Time{}, fmt.Errorf("failed to convert time.Time from bytes", bv)
}

func (bv BytesValue) ToRat() (*big.Rat, error) {
	r := new(big.Rat)
	r.SetString(string(bv))
	return r, nil
}

func (bv BytesValue) Marshal() (string, error) {
	v, err := bv.ToString()
	if err != nil {
		return "", err
	}
	return strconv.Quote(v), nil
}

func (bv BytesValue) Format(verb rune) string {
	v, _ := bv.ToString()
	switch verb {
	case 't':
		return v
	case 'T':
		return strconv.Quote(v)
	}
	return v
}

func (bv BytesValue) Interface() interface{} {
	return []byte(bv)
}

type FloatValue float64

func (fv FloatValue) Add(v Value) (Value, error) {
	v2, err := v.ToFloat64()
	if err != nil {
		return nil, err
	}
	return ValueOf(float64(fv) + v2)
}

func (fv FloatValue) Sub(v Value) (Value, error) {
	v2, err := v.ToFloat64()
	if err != nil {
		return nil, err
	}
	return ValueOf(float64(fv) - v2)
}

func (fv FloatValue) Mul(v Value) (Value, error) {
	v2, err := v.ToFloat64()
	if err != nil {
		return nil, err
	}
	return ValueOf(float64(fv) * v2)
}

func (fv FloatValue) Div(v Value) (Value, error) {
	v2, err := v.ToFloat64()
	if err != nil {
		return nil, err
	}
	if v2 == 0 {
		return nil, fmt.Errorf("zero divided error ( %f / 0 )", fv)
	}
	return ValueOf(float64(fv) / v2)
}

func (fv FloatValue) EQ(v Value) (bool, error) {
	v2, err := v.ToFloat64()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to float64", v)
	}
	return float64(fv) == v2, nil
}

func (fv FloatValue) GT(v Value) (bool, error) {
	v2, err := v.ToFloat64()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to float64", v)
	}
	return float64(fv) > v2, nil
}

func (fv FloatValue) GTE(v Value) (bool, error) {
	v2, err := v.ToFloat64()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to float64", v)
	}
	return float64(fv) >= v2, nil
}

func (fv FloatValue) LT(v Value) (bool, error) {
	v2, err := v.ToFloat64()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to float64", v)
	}
	return float64(fv) < v2, nil
}

func (fv FloatValue) LTE(v Value) (bool, error) {
	v2, err := v.ToFloat64()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to float64", v)
	}
	return float64(fv) <= v2, nil
}

func (fv FloatValue) ToInt64() (int64, error) {
	return int64(fv), nil
}

func (fv FloatValue) ToString() (string, error) {
	return fmt.Sprint(fv), nil
}

func (fv FloatValue) ToBytes() ([]byte, error) {
	return []byte(fmt.Sprint(fv)), nil
}

func (fv FloatValue) ToFloat64() (float64, error) {
	return float64(fv), nil
}

func (fv FloatValue) ToBool() (bool, error) {
	return false, fmt.Errorf("failed to convert %f to bool type", fv)
}

func (fv FloatValue) ToArray() (*ArrayValue, error) {
	return nil, fmt.Errorf("failed to convert array from float64: %v", fv)
}

func (fv FloatValue) ToStruct() (*StructValue, error) {
	return nil, fmt.Errorf("failed to convert struct from float64: %v", fv)
}

func (fv FloatValue) ToJSON() (string, error) {
	return fmt.Sprint(fv), nil
}

func (fv FloatValue) ToTime() (time.Time, error) {
	return time.Time{}, fmt.Errorf("failed to convert time.Time from float64: %v", fv)
}

func (fv FloatValue) ToRat() (*big.Rat, error) {
	r := new(big.Rat)
	r.SetFloat64(float64(fv))
	return r, nil
}

func (fv FloatValue) Marshal() (string, error) {
	return fmt.Sprint(fv), nil
}

func (fv FloatValue) Format(verb rune) string {
	return fmt.Sprint(fv)
}

func (fv FloatValue) Interface() interface{} {
	return float64(fv)
}

type NumericValue big.Rat

func (nv *NumericValue) Add(v Value) (Value, error) {
	z := new(big.Rat)
	x := (*big.Rat)(nv)
	y, err := v.ToRat()
	if err != nil {
		return nil, err
	}
	return (*NumericValue)(z.Add(x, y)), nil
}

func (nv *NumericValue) Sub(v Value) (Value, error) {
	z := new(big.Rat)
	x := (*big.Rat)(nv)
	y, err := v.ToRat()
	if err != nil {
		return nil, err
	}
	zy := new(big.Rat)
	return (*NumericValue)(z.Add(x, zy.Neg(y))), nil
}

func (nv *NumericValue) Mul(v Value) (Value, error) {
	z := new(big.Rat)
	x := (*big.Rat)(nv)
	y, err := v.ToRat()
	if err != nil {
		return nil, err
	}
	return (*NumericValue)(z.Mul(x, y)), nil
}

func (nv *NumericValue) Div(v Value) (ret Value, e error) {
	defer func() {
		if err := recover(); err != nil {
			e = err.(error)
		}
	}()
	z := new(big.Rat)
	x := (*big.Rat)(nv)
	y, err := v.ToRat()
	if err != nil {
		return nil, err
	}
	zy := new(big.Rat)
	return (*NumericValue)(z.Mul(x, zy.Inv(y))), nil
}

func (nv *NumericValue) EQ(v Value) (bool, error) {
	x := (*big.Rat)(nv)
	y, err := v.ToRat()
	if err != nil {
		return false, err
	}
	return x.Cmp(y) == 0, nil
}

func (nv *NumericValue) GT(v Value) (bool, error) {
	x := (*big.Rat)(nv)
	y, err := v.ToRat()
	if err != nil {
		return false, err
	}
	return x.Cmp(y) > 0, nil
}

func (nv *NumericValue) GTE(v Value) (bool, error) {
	x := (*big.Rat)(nv)
	y, err := v.ToRat()
	if err != nil {
		return false, err
	}
	return x.Cmp(y) >= 0, nil
}

func (nv *NumericValue) LT(v Value) (bool, error) {
	x := (*big.Rat)(nv)
	y, err := v.ToRat()
	if err != nil {
		return false, err
	}
	return x.Cmp(y) < 0, nil
}

func (nv *NumericValue) LTE(v Value) (bool, error) {
	x := (*big.Rat)(nv)
	y, err := v.ToRat()
	if err != nil {
		return false, err
	}
	return x.Cmp(y) <= 0, nil
}

func (nv *NumericValue) ToInt64() (int64, error) {
	return (*big.Rat)(nv).Num().Int64(), nil
}

func (nv *NumericValue) ToString() (string, error) {
	return (*big.Rat)(nv).RatString(), nil
}

func (nv *NumericValue) ToBytes() ([]byte, error) {
	return []byte((*big.Rat)(nv).RatString()), nil
}

func (nv *NumericValue) ToFloat64() (float64, error) {
	f, _ := (*big.Rat)(nv).Float64()
	return f, nil
}

func (nv *NumericValue) ToBool() (bool, error) {
	v := (*big.Rat)(nv).Num().Int64()
	if v == 1 {
		return true, nil
	} else if v == 0 {
		return false, nil
	}
	return false, fmt.Errorf("failed to convert numeric value to bool type")
}

func (nv *NumericValue) ToArray() (*ArrayValue, error) {
	return nil, fmt.Errorf("failed to convert array from numeric value")
}

func (nv *NumericValue) ToStruct() (*StructValue, error) {
	return nil, fmt.Errorf("failed to convert struct from numeric value")
}

func (nv *NumericValue) ToJSON() (string, error) {
	return (*big.Rat)(nv).RatString(), nil
}

func (nv *NumericValue) ToTime() (time.Time, error) {
	return time.Time{}, fmt.Errorf("failed to convert time.Time from numeric value")
}

func (nv *NumericValue) ToRat() (*big.Rat, error) {
	return (*big.Rat)(nv), nil
}

func (nv *NumericValue) Marshal() (string, error) {
	b, err := (*big.Rat)(nv).MarshalText()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (nv *NumericValue) Format(verb rune) string {
	return (*big.Rat)(nv).RatString()
}

func (nv *NumericValue) Interface() interface{} {
	f, _ := (*big.Rat)(nv).Float64()
	return f
}

type BoolValue bool

func (bv BoolValue) Add(v Value) (Value, error) {
	return nil, fmt.Errorf("add operation is unsupported for bool %v", bv)
}

func (bv BoolValue) Sub(v Value) (Value, error) {
	return nil, fmt.Errorf("sub operation is unsupported for bool %v", bv)
}

func (bv BoolValue) Mul(v Value) (Value, error) {
	return nil, fmt.Errorf("mul operation is unsupported for bool %v", bv)
}

func (bv BoolValue) Div(v Value) (Value, error) {
	return nil, fmt.Errorf("div operation is unsupported for bool %v", bv)
}

func (bv BoolValue) EQ(v Value) (bool, error) {
	v2, err := v.ToBool()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to bool", v)
	}
	return bool(bv) == v2, nil
}

func (bv BoolValue) GT(v Value) (bool, error) {
	return false, fmt.Errorf("gt operation is unsupported for bool %v", bv)
}

func (bv BoolValue) GTE(v Value) (bool, error) {
	return false, fmt.Errorf("gte operation is unsupported for bool %v", bv)
}

func (bv BoolValue) LT(v Value) (bool, error) {
	return false, fmt.Errorf("lt operation is unsupported for bool %v", bv)
}

func (bv BoolValue) LTE(v Value) (bool, error) {
	return false, fmt.Errorf("lte operation is unsupported for bool %v", bv)
}

func (bv BoolValue) ToInt64() (int64, error) {
	if bv {
		return 1, nil
	}
	return 0, nil
}

func (bv BoolValue) ToString() (string, error) {
	return fmt.Sprint(bv), nil
}

func (bv BoolValue) ToBytes() ([]byte, error) {
	return []byte(fmt.Sprint(bv)), nil
}

func (bv BoolValue) ToFloat64() (float64, error) {
	if bv {
		return 1, nil
	}
	return 0, nil
}

func (bv BoolValue) ToBool() (bool, error) {
	return bool(bv), nil
}

func (bv BoolValue) ToArray() (*ArrayValue, error) {
	return nil, fmt.Errorf("failed to convert bool from array: %v", bv)
}

func (bv BoolValue) ToStruct() (*StructValue, error) {
	return nil, fmt.Errorf("failed to convert bool from struct: %v", bv)
}

func (bv BoolValue) ToJSON() (string, error) {
	return fmt.Sprint(bv), nil
}

func (bv BoolValue) ToTime() (time.Time, error) {
	return time.Time{}, fmt.Errorf("failed to convert bool from time.Time: %v", bv)
}

func (bv BoolValue) ToRat() (*big.Rat, error) {
	r := new(big.Rat)
	if bv {
		r.SetInt64(1)
		return r, nil
	}
	r.SetInt64(0)
	return r, nil
}

func (bv BoolValue) Marshal() (string, error) {
	return fmt.Sprint(bv), nil
}

func (bv BoolValue) Format(verb rune) string {
	return fmt.Sprint(bv)
}

func (bv BoolValue) Interface() interface{} {
	return bool(bv)
}

type JsonValue string

func (jv JsonValue) Add(v Value) (Value, error) {
	return nil, fmt.Errorf("add operation is unsupported for json %v", jv)
}

func (jv JsonValue) Sub(v Value) (Value, error) {
	return nil, fmt.Errorf("sub operation is unsupported for json %v", jv)
}

func (jv JsonValue) Mul(v Value) (Value, error) {
	return nil, fmt.Errorf("mul operation is unsupported for json %v", jv)
}

func (jv JsonValue) Div(v Value) (Value, error) {
	return nil, fmt.Errorf("div operation is unsupported for json %v", jv)
}

func (jv JsonValue) EQ(v Value) (bool, error) {
	return false, fmt.Errorf("eq operation is unsupported for json %v", jv)
}

func (jv JsonValue) GT(v Value) (bool, error) {
	return false, fmt.Errorf("gt operation is unsupported for json %v", jv)
}

func (jv JsonValue) GTE(v Value) (bool, error) {
	return false, fmt.Errorf("gte operation is unsupported for json %v", jv)
}

func (jv JsonValue) LT(v Value) (bool, error) {
	return false, fmt.Errorf("lt operation is unsupported for json %v", jv)
}

func (jv JsonValue) LTE(v Value) (bool, error) {
	return false, fmt.Errorf("lte operation is unsupported for json %v", jv)
}

func (jv JsonValue) ToInt64() (int64, error) {
	return strconv.ParseInt(string(jv), 10, 64)
}

func (jv JsonValue) ToString() (string, error) {
	return toJsonValueFromString(string(jv))
}

func (jv JsonValue) ToBytes() ([]byte, error) {
	v, err := toJsonValueFromString(string(jv))
	if err != nil {
		return nil, err
	}
	return []byte(v), nil
}

func (jv JsonValue) ToFloat64() (float64, error) {
	return strconv.ParseFloat(string(jv), 64)
}

func (jv JsonValue) ToBool() (bool, error) {
	return strconv.ParseBool(string(jv))
}

func (jv JsonValue) ToArray() (*ArrayValue, error) {
	return nil, fmt.Errorf("failed to convert json from array: %v", jv)
}

func (jv JsonValue) ToStruct() (*StructValue, error) {
	return nil, fmt.Errorf("failed to convert json from struct: %v", jv)
}

func (jv JsonValue) ToJSON() (string, error) {
	return string(jv), nil
}

func (jv JsonValue) ToTime() (time.Time, error) {
	return time.Time{}, fmt.Errorf("failed to convert json from time.Time: %v", jv)
}

func (jv JsonValue) ToRat() (*big.Rat, error) {
	i64, err := strconv.ParseInt(string(jv), 10, 64)
	if err != nil {
		return nil, err
	}
	r := new(big.Rat)
	r.SetInt64(i64)
	return r, nil
}

func (jv JsonValue) Marshal() (string, error) {
	return jv.ToString()
}

func (jv JsonValue) Format(verb rune) string {
	return string(jv)
}

func (jv JsonValue) Interface() interface{} {
	var v interface{}
	if err := json.Unmarshal([]byte(jv), &v); err != nil {
		return nil
	}
	return v
}

func (jv JsonValue) reflectTypeToJsonType(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return "number"
	case reflect.String:
		return "string"
	case reflect.Bool:
		return "boolean"
	case reflect.Slice, reflect.Array:
		return "array"
	case reflect.Struct, reflect.Map:
		return "object"
	case reflect.Ptr:
		return jv.reflectTypeToJsonType(t.Elem())
	}
	return "unknown"
}

func (jv JsonValue) Type() string {
	if string(jv) == "null" {
		return "null"
	}
	rv := reflect.ValueOf(jv.Interface())
	return jv.reflectTypeToJsonType(rv.Type())
}

type ArrayValue struct {
	values []Value
}

func (av *ArrayValue) Has(v Value) (bool, error) {
	for _, val := range av.values {
		cond, err := val.EQ(v)
		if err != nil {
			return false, err
		}
		if cond {
			return true, nil
		}
	}
	return false, nil
}

func (av *ArrayValue) Add(v Value) (Value, error) {
	return nil, fmt.Errorf("add operation is unsupported for array %v", av)
}

func (av *ArrayValue) Sub(v Value) (Value, error) {
	return nil, fmt.Errorf("sub operation is unsupported for array %v", av)
}

func (av *ArrayValue) Mul(v Value) (Value, error) {
	return nil, fmt.Errorf("mul operation is unsupported for array %v", av)
}

func (av *ArrayValue) Div(v Value) (Value, error) {
	return nil, fmt.Errorf("div operation is unsupported for array %v", av)
}

func (av *ArrayValue) EQ(v Value) (bool, error) {
	arr, err := v.ToArray()
	if err != nil {
		return false, err
	}
	if len(arr.values) != len(av.values) {
		return false, nil
	}
	for idx, value := range av.values {
		cond, err := arr.values[idx].EQ(value)
		if err != nil {
			return false, err
		}
		if !cond {
			return false, nil
		}
	}
	return true, nil
}

func (av *ArrayValue) GT(v Value) (bool, error) {
	arr, err := v.ToArray()
	if err != nil {
		return false, err
	}
	if len(arr.values) != len(av.values) {
		return false, nil
	}
	for idx, value := range av.values {
		cond, err := arr.values[idx].GT(value)
		if err != nil {
			return false, err
		}
		if !cond {
			return false, nil
		}
	}
	return true, nil
}

func (av *ArrayValue) GTE(v Value) (bool, error) {
	arr, err := v.ToArray()
	if err != nil {
		return false, err
	}
	if len(arr.values) != len(av.values) {
		return false, nil
	}
	for idx, value := range av.values {
		cond, err := arr.values[idx].GTE(value)
		if err != nil {
			return false, err
		}
		if !cond {
			return false, nil
		}
	}
	return true, nil
}

func (av *ArrayValue) LT(v Value) (bool, error) {
	arr, err := v.ToArray()
	if err != nil {
		return false, err
	}
	if len(arr.values) != len(av.values) {
		return false, nil
	}
	for idx, value := range av.values {
		cond, err := arr.values[idx].LT(value)
		if err != nil {
			return false, err
		}
		if !cond {
			return false, nil
		}
	}
	return true, nil
}

func (av *ArrayValue) LTE(v Value) (bool, error) {
	arr, err := v.ToArray()
	if err != nil {
		return false, err
	}
	if len(arr.values) != len(av.values) {
		return false, nil
	}
	for idx, value := range av.values {
		cond, err := arr.values[idx].LTE(value)
		if err != nil {
			return false, err
		}
		if !cond {
			return false, nil
		}
	}
	return true, nil
}

func (av *ArrayValue) ToInt64() (int64, error) {
	return 0, fmt.Errorf("failed to convert int64 from array %v", av)
}

func (av *ArrayValue) ToString() (string, error) {
	return av.Marshal()
}

func (av *ArrayValue) ToBytes() ([]byte, error) {
	v, err := av.Marshal()
	if err != nil {
		return nil, err
	}
	return []byte(v), nil
}

func (av *ArrayValue) ToFloat64() (float64, error) {
	return 0, fmt.Errorf("failed to convert float64 from array %v", av)
}

func (av *ArrayValue) ToBool() (bool, error) {
	return false, fmt.Errorf("failed to convert bool from array %v", av)
}

func (av *ArrayValue) ToArray() (*ArrayValue, error) {
	return av, nil
}

func (av *ArrayValue) ToStruct() (*StructValue, error) {
	return nil, fmt.Errorf("failed to convert struct from array %v", av)
}

func (av *ArrayValue) ToJSON() (string, error) {
	elems := []string{}
	for _, v := range av.values {
		if v == nil {
			elems = append(elems, "null")
			continue
		}
		elem, err := v.ToJSON()
		if err != nil {
			return "", err
		}
		elems = append(elems, elem)
	}
	return fmt.Sprintf("[%s]", strings.Join(elems, ",")), nil
}

func (av *ArrayValue) ToTime() (time.Time, error) {
	return time.Time{}, fmt.Errorf("failed to convert time.Time from array %v", av)
}

func (av *ArrayValue) ToRat() (*big.Rat, error) {
	return nil, fmt.Errorf("failed to convert *big.Rat from array %v", av)
}

func (av *ArrayValue) Marshal() (string, error) {
	elems := []string{}
	for _, v := range av.values {
		if v == nil {
			elems = append(elems, "null")
			continue
		}
		elem, err := v.Marshal()
		if err != nil {
			return "", err
		}
		elems = append(elems, elem)
	}
	return toArrayValueFromJSONString(fmt.Sprintf("[%s]", strings.Join(elems, ","))), nil
}

func (av *ArrayValue) Format(verb rune) string {
	elems := []string{}
	for _, v := range av.values {
		if v == nil {
			elems = append(elems, "NULL")
			continue
		}
		elems = append(elems, v.Format(verb))
	}
	return fmt.Sprintf("[%s]", strings.Join(elems, ", "))
}

func (av *ArrayValue) Interface() interface{} {
	var arr []interface{}
	for _, v := range av.values {
		if v == nil {
			arr = append(arr, nil)
		} else {
			arr = append(arr, v.Interface())
		}
	}
	return arr
}

type StructValue struct {
	keys   []string
	values []Value
	m      map[string]Value
}

func (sv *StructValue) Add(v Value) (Value, error) {
	return nil, fmt.Errorf("add operation is unsupported for struct %v", sv)
}

func (sv *StructValue) Sub(v Value) (Value, error) {
	return nil, fmt.Errorf("sub operation is unsupported for struct %v", sv)
}

func (sv *StructValue) Mul(v Value) (Value, error) {
	return nil, fmt.Errorf("mul operation is unsupported for struct %v", sv)
}

func (sv *StructValue) Div(v Value) (Value, error) {
	return nil, fmt.Errorf("div operation is unsupported for struct %v", sv)
}

func (sv *StructValue) EQ(v Value) (bool, error) {
	st, err := v.ToStruct()
	if err != nil {
		return false, err
	}
	if len(st.m) != len(sv.m) {
		return false, nil
	}
	for key := range sv.m {
		cond, err := st.m[key].EQ(sv.m[key])
		if err != nil {
			return false, err
		}
		if !cond {
			return false, nil
		}
	}
	return true, nil
}

func (sv *StructValue) GT(v Value) (bool, error) {
	st, err := v.ToStruct()
	if err != nil {
		return false, err
	}
	if len(st.m) != len(sv.m) {
		return false, nil
	}
	for key := range sv.m {
		cond, err := st.m[key].GT(sv.m[key])
		if err != nil {
			return false, err
		}
		if !cond {
			return false, nil
		}
	}
	return true, nil
}

func (sv *StructValue) GTE(v Value) (bool, error) {
	st, err := v.ToStruct()
	if err != nil {
		return false, err
	}
	if len(st.m) != len(sv.m) {
		return false, nil
	}
	for key := range sv.m {
		cond, err := st.m[key].GTE(sv.m[key])
		if err != nil {
			return false, err
		}
		if !cond {
			return false, nil
		}
	}
	return true, nil
}

func (sv *StructValue) LT(v Value) (bool, error) {
	st, err := v.ToStruct()
	if err != nil {
		return false, err
	}
	if len(st.m) != len(sv.m) {
		return false, nil
	}
	for key := range sv.m {
		cond, err := st.m[key].LT(sv.m[key])
		if err != nil {
			return false, err
		}
		if !cond {
			return false, nil
		}
	}
	return true, nil
}

func (sv *StructValue) LTE(v Value) (bool, error) {
	st, err := v.ToStruct()
	if err != nil {
		return false, err
	}
	if len(st.m) != len(sv.m) {
		return false, nil
	}
	for key := range sv.m {
		cond, err := st.m[key].LTE(sv.m[key])
		if err != nil {
			return false, err
		}
		if !cond {
			return false, nil
		}
	}
	return true, nil
}

func (sv *StructValue) ToInt64() (int64, error) {
	return 0, fmt.Errorf("failed to convert int64 from struct %v", sv)
}

func (sv *StructValue) ToString() (string, error) {
	return sv.Marshal()
}

func (sv *StructValue) ToBytes() ([]byte, error) {
	v, err := sv.Marshal()
	if err != nil {
		return nil, err
	}
	return []byte(v), nil
}

func (sv *StructValue) ToFloat64() (float64, error) {
	return 0, fmt.Errorf("failed to convert float64 from struct %v", sv)
}

func (sv *StructValue) ToBool() (bool, error) {
	return false, fmt.Errorf("failed to convert bool from struct %v", sv)
}

func (sv *StructValue) ToArray() (*ArrayValue, error) {
	return nil, fmt.Errorf("failed to convert array from struct %v", sv)
}

func (sv *StructValue) ToStruct() (*StructValue, error) {
	return sv, nil
}

func (sv *StructValue) ToJSON() (string, error) {
	fields := []string{}
	for i := 0; i < len(sv.keys); i++ {
		key := sv.keys[i]
		value, err := sv.values[i].ToJSON()
		if err != nil {
			return "", err
		}
		fields = append(
			fields,
			fmt.Sprintf("%s:%s", strconv.Quote(key), value),
		)
	}
	return fmt.Sprintf("{%s}", strings.Join(fields, ",")), nil
}

func (sv *StructValue) ToTime() (time.Time, error) {
	return time.Time{}, fmt.Errorf("failed to convert time.Time from struct %v", sv)
}

func (sv *StructValue) ToRat() (*big.Rat, error) {
	return nil, fmt.Errorf("failed to convert *big.Rat from struct %v", sv)
}

func (sv *StructValue) Marshal() (string, error) {
	fields := []string{}
	for i := 0; i < len(sv.keys); i++ {
		key := sv.keys[i]
		value := sv.values[i]
		if value == nil {
			fields = append(
				fields,
				fmt.Sprintf("%s:null", strconv.Quote(key)),
			)
			continue
		}
		encodedValue, err := value.Marshal()
		if err != nil {
			return "", err
		}
		fields = append(
			fields,
			fmt.Sprintf("%s:%s", strconv.Quote(key), encodedValue),
		)
	}
	return toStructValueFromJSONString(
		fmt.Sprintf("{%s}", strings.Join(fields, ",")),
	), nil
}

func (sv *StructValue) Format(verb rune) string {
	elems := []string{}
	for _, v := range sv.values {
		if v == nil {
			elems = append(elems, "NULL")
			continue
		}
		elems = append(elems, v.Format(verb))
	}
	return fmt.Sprintf("(%s)", strings.Join(elems, ", "))
}

func (sv *StructValue) Interface() interface{} {
	fields := []map[string]interface{}{}
	for i := 0; i < len(sv.keys); i++ {
		fields = append(fields, map[string]interface{}{
			sv.keys[i]: sv.values[i].Interface(),
		})
	}
	return fields
}

type DateValue time.Time

func (d DateValue) AddDateWithInterval(v int, interval string) (Value, error) {
	switch interval {
	case "WEEK":
		return DateValue(time.Time(d).AddDate(0, 0, v*7)), nil
	case "MONTH":
		return DateValue(time.Time(d).AddDate(0, v, 0)), nil
	case "YEAR":
		return DateValue(time.Time(d).AddDate(v, 0, 0)), nil
	default:
		return DateValue(time.Time(d).AddDate(0, 0, v)), nil
	}
}

func (d DateValue) Add(v Value) (Value, error) {
	v2, err := v.ToInt64()
	if err != nil {
		return nil, err
	}
	duration := time.Duration(v2) * 24 * time.Hour
	return DateValue(time.Time(d).Add(duration)), nil
}

func (d DateValue) Sub(v Value) (Value, error) {
	v2, err := v.ToInt64()
	if err != nil {
		return nil, err
	}
	duration := -time.Duration(v2) * 24 * time.Hour
	return DateValue(time.Time(d).Add(duration)), nil
}

func (d DateValue) Mul(v Value) (Value, error) {
	return nil, fmt.Errorf("mul operation is unsupported for date %v", d)
}

func (d DateValue) Div(v Value) (Value, error) {
	return nil, fmt.Errorf("div operation is unsupported for date %v", d)
}

func (d DateValue) EQ(v Value) (bool, error) {
	v2, err := v.ToTime()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to time.Time", v)
	}
	return time.Time(d).Equal(v2), nil
}

func (d DateValue) GT(v Value) (bool, error) {
	v2, err := v.ToTime()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to time.Time", v)
	}
	return time.Time(d).After(v2), nil
}

func (d DateValue) GTE(v Value) (bool, error) {
	v2, err := v.ToTime()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to time.Time", v)
	}
	return time.Time(d).Equal(v2) || time.Time(d).After(v2), nil
}

func (d DateValue) LT(v Value) (bool, error) {
	v2, err := v.ToTime()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to time.Time", v)
	}
	return time.Time(d).Before(v2), nil
}

func (d DateValue) LTE(v Value) (bool, error) {
	v2, err := v.ToTime()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to time.Time", v)
	}
	return time.Time(d).Equal(v2) || time.Time(d).Before(v2), nil
}

func (d DateValue) ToInt64() (int64, error) {
	return time.Time(d).Unix(), nil
}

func (d DateValue) ToString() (string, error) {
	json, err := d.ToJSON()
	if err != nil {
		return "", err
	}
	return toDateValueFromString(json), nil
}

func (d DateValue) ToBytes() ([]byte, error) {
	json, err := d.ToJSON()
	if err != nil {
		return nil, err
	}
	return []byte(toDateValueFromString(json)), nil
}

func (d DateValue) ToFloat64() (float64, error) {
	return float64(time.Time(d).Unix()), nil
}

func (d DateValue) ToBool() (bool, error) {
	return false, fmt.Errorf("failed to convert %v to bool type", d)
}

func (d DateValue) ToArray() (*ArrayValue, error) {
	return nil, fmt.Errorf("failed to convert %v to array type", d)
}

func (d DateValue) ToStruct() (*StructValue, error) {
	return nil, fmt.Errorf("failed to convert %v to struct type", d)
}

func (d DateValue) ToJSON() (string, error) {
	return time.Time(d).Format("2006-01-02"), nil
}

func (d DateValue) ToTime() (time.Time, error) {
	return time.Time(d), nil
}

func (d DateValue) ToRat() (*big.Rat, error) {
	return nil, fmt.Errorf("failed to convert *big.Rat from date %v", d)
}

func (d DateValue) Marshal() (string, error) {
	json, err := d.ToJSON()
	if err != nil {
		return "", err
	}
	return toDateValueFromString(json), nil
}

func (d DateValue) Format(verb rune) string {
	formatted := time.Time(d).Format("2006-01-02")
	switch verb {
	case 't':
		return formatted
	case 'T':
		return fmt.Sprintf(`DATE "%s"`, formatted)
	}
	return formatted
}

func (d DateValue) Interface() interface{} {
	return time.Time(d).Format("2006-01-02")
}

type DatetimeValue time.Time

func (d DatetimeValue) Add(v Value) (Value, error) {
	v2, err := v.ToInt64()
	if err != nil {
		return nil, err
	}
	duration := time.Duration(v2) * 24 * time.Hour
	return DatetimeValue(time.Time(d).Add(duration)), nil
}

func (d DatetimeValue) Sub(v Value) (Value, error) {
	v2, err := v.ToInt64()
	if err != nil {
		return nil, err
	}
	duration := -time.Duration(v2) * 24 * time.Hour
	return DateValue(time.Time(d).Add(duration)), nil
}

func (d DatetimeValue) Mul(v Value) (Value, error) {
	return nil, fmt.Errorf("mul operation is unsupported for datetime %v", d)
}

func (d DatetimeValue) Div(v Value) (Value, error) {
	return nil, fmt.Errorf("div operation is unsupported for datetime %v", d)
}

func (d DatetimeValue) EQ(v Value) (bool, error) {
	v2, err := v.ToTime()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to time.Time", v)
	}
	return time.Time(d).Equal(v2), nil
}

func (d DatetimeValue) GT(v Value) (bool, error) {
	v2, err := v.ToTime()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to time.Time", v)
	}
	return time.Time(d).After(v2), nil
}

func (d DatetimeValue) GTE(v Value) (bool, error) {
	v2, err := v.ToTime()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to time.Time", v)
	}
	return time.Time(d).Equal(v2) || time.Time(d).After(v2), nil
}

func (d DatetimeValue) LT(v Value) (bool, error) {
	v2, err := v.ToTime()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to time.Time", v)
	}
	return time.Time(d).Before(v2), nil
}

func (d DatetimeValue) LTE(v Value) (bool, error) {
	v2, err := v.ToTime()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to time.Time", v)
	}
	return time.Time(d).Equal(v2) || time.Time(d).Before(v2), nil
}

func (d DatetimeValue) ToInt64() (int64, error) {
	return time.Time(d).Unix(), nil
}

func (d DatetimeValue) ToString() (string, error) {
	json, err := d.ToJSON()
	if err != nil {
		return "", err
	}
	return toDatetimeValueFromString(json), nil
}

func (d DatetimeValue) ToBytes() ([]byte, error) {
	json, err := d.ToJSON()
	if err != nil {
		return nil, err
	}
	return []byte(toDatetimeValueFromString(json)), nil
}

func (d DatetimeValue) ToFloat64() (float64, error) {
	return float64(time.Time(d).Unix()), nil
}

func (d DatetimeValue) ToBool() (bool, error) {
	return false, fmt.Errorf("failed to convert %v to bool type", d)
}

func (d DatetimeValue) ToArray() (*ArrayValue, error) {
	return nil, fmt.Errorf("failed to convert %v to array type", d)
}

func (d DatetimeValue) ToStruct() (*StructValue, error) {
	return nil, fmt.Errorf("failed to convert %v to struct type", d)
}

func (d DatetimeValue) ToJSON() (string, error) {
	return time.Time(d).Format("2006-01-02T15:04:05"), nil
}

func (d DatetimeValue) ToTime() (time.Time, error) {
	return time.Time(d), nil
}

func (d DatetimeValue) ToRat() (*big.Rat, error) {
	return nil, fmt.Errorf("failed to convert *big.Rat from datetime %v", d)
}

func (d DatetimeValue) Marshal() (string, error) {
	json, err := d.ToJSON()
	if err != nil {
		return "", err
	}
	return toDatetimeValueFromString(json), nil
}

func (d DatetimeValue) Format(verb rune) string {
	formatted := time.Time(d).Format("2006-01-02T15:04:05")
	switch verb {
	case 't':
		return formatted
	case 'T':
		return fmt.Sprintf(`DATETIME "%s"`, formatted)
	}
	return formatted
}

func (d DatetimeValue) Interface() interface{} {
	return time.Time(d).Format("2006-01-02T15:04:05")
}

type TimeValue time.Time

func (d TimeValue) Add(v Value) (Value, error) {
	v2, err := v.ToInt64()
	if err != nil {
		return nil, err
	}
	duration := time.Duration(v2) * 24 * time.Hour
	return TimeValue(time.Time(d).Add(duration)), nil
}

func (d TimeValue) Sub(v Value) (Value, error) {
	v2, err := v.ToInt64()
	if err != nil {
		return nil, err
	}
	duration := -time.Duration(v2) * 24 * time.Hour
	return TimeValue(time.Time(d).Add(duration)), nil
}

func (d TimeValue) Mul(v Value) (Value, error) {
	return nil, fmt.Errorf("mul operation is unsupported for time %v", d)
}

func (d TimeValue) Div(v Value) (Value, error) {
	return nil, fmt.Errorf("div operation is unsupported for time %v", d)
}

func (d TimeValue) EQ(v Value) (bool, error) {
	v2, err := v.ToTime()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to time.Time", v)
	}
	return time.Time(d).Equal(v2), nil
}

func (d TimeValue) GT(v Value) (bool, error) {
	v2, err := v.ToTime()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to time.Time", v)
	}
	return time.Time(d).After(v2), nil
}

func (d TimeValue) GTE(v Value) (bool, error) {
	v2, err := v.ToTime()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to time.Time", v)
	}
	return time.Time(d).Equal(v2) || time.Time(d).After(v2), nil
}

func (d TimeValue) LT(v Value) (bool, error) {
	v2, err := v.ToTime()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to time.Time", v)
	}
	return time.Time(d).Before(v2), nil
}

func (d TimeValue) LTE(v Value) (bool, error) {
	v2, err := v.ToTime()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to time.Time", v)
	}
	return time.Time(d).Equal(v2) || time.Time(d).Before(v2), nil
}

func (d TimeValue) ToInt64() (int64, error) {
	return time.Time(d).Unix(), nil
}

func (d TimeValue) ToString() (string, error) {
	json, err := d.ToJSON()
	if err != nil {
		return "", err
	}
	return toTimeValueFromString(json), nil
}

func (d TimeValue) ToBytes() ([]byte, error) {
	json, err := d.ToJSON()
	if err != nil {
		return nil, err
	}
	return []byte(toTimeValueFromString(json)), nil
}

func (d TimeValue) ToFloat64() (float64, error) {
	return float64(time.Time(d).Unix()), nil
}

func (d TimeValue) ToBool() (bool, error) {
	return false, fmt.Errorf("failed to convert %v to bool type", d)
}

func (d TimeValue) ToArray() (*ArrayValue, error) {
	return nil, fmt.Errorf("failed to convert %v to array type", d)
}

func (d TimeValue) ToStruct() (*StructValue, error) {
	return nil, fmt.Errorf("failed to convert %v to struct type", d)
}

func (d TimeValue) ToJSON() (string, error) {
	return time.Time(d).Format("15:04:05"), nil
}

func (d TimeValue) ToTime() (time.Time, error) {
	return time.Time(d), nil
}

func (d TimeValue) ToRat() (*big.Rat, error) {
	return nil, fmt.Errorf("failed to convert *big.Rat from time %v", d)
}

func (d TimeValue) Marshal() (string, error) {
	json, err := d.ToJSON()
	if err != nil {
		return "", err
	}
	return toTimeValueFromString(json), nil
}

func (d TimeValue) Format(verb rune) string {
	formatted := time.Time(d).Format("15:04:05")
	switch verb {
	case 't':
		return formatted
	case 'T':
		return fmt.Sprintf(`TIME "%s"`, formatted)
	}
	return formatted
}

func (d TimeValue) Interface() interface{} {
	return time.Time(d).Format("15:04:05")
}

type TimestampValue time.Time

func (d TimestampValue) AddValueWithPart(v time.Duration, part string) (Value, error) {
	switch part {
	case "MICROSECOND":
		return TimestampValue(time.Time(d).Add(v * time.Microsecond)), nil
	case "MILLISECOND":
		return TimestampValue(time.Time(d).Add(v * time.Millisecond)), nil
	case "SECOND":
		return TimestampValue(time.Time(d).Add(v * time.Second)), nil
	case "MINUTE":
		return TimestampValue(time.Time(d).Add(v * time.Minute)), nil
	case "HOUR":
		return TimestampValue(time.Time(d).Add(v * time.Hour)), nil
	case "DAY":
		return TimestampValue(time.Time(d).Add(v * time.Hour * 24)), nil
	default:
		return nil, fmt.Errorf("unknown part value for timestamp: %s", part)
	}
}

func (d TimestampValue) Add(v Value) (Value, error) {
	v2, err := v.ToInt64()
	if err != nil {
		return nil, err
	}
	duration := time.Duration(v2) * 24 * time.Hour
	return TimestampValue(time.Time(d).Add(duration)), nil
}

func (d TimestampValue) Sub(v Value) (Value, error) {
	v2, err := v.ToInt64()
	if err != nil {
		return nil, err
	}
	duration := -time.Duration(v2) * 24 * time.Hour
	return TimestampValue(time.Time(d).Add(duration)), nil
}

func (d TimestampValue) Mul(v Value) (Value, error) {
	return nil, fmt.Errorf("mul operation is unsupported for timestamp %v", d)
}

func (d TimestampValue) Div(v Value) (Value, error) {
	return nil, fmt.Errorf("div operation is unsupported for timestamp %v", d)
}

func (d TimestampValue) EQ(v Value) (bool, error) {
	v2, err := v.ToTime()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to time.Time", v)
	}
	return time.Time(d).Equal(v2), nil
}

func (d TimestampValue) GT(v Value) (bool, error) {
	v2, err := v.ToTime()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to time.Time", v)
	}
	return time.Time(d).After(v2), nil
}

func (d TimestampValue) GTE(v Value) (bool, error) {
	v2, err := v.ToTime()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to time.Time", v)
	}
	return time.Time(d).Equal(v2) || time.Time(d).After(v2), nil
}

func (d TimestampValue) LT(v Value) (bool, error) {
	v2, err := v.ToTime()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to time.Time", v)
	}
	return time.Time(d).Before(v2), nil
}

func (d TimestampValue) LTE(v Value) (bool, error) {
	v2, err := v.ToTime()
	if err != nil {
		return false, fmt.Errorf("failed to convert %v to time.Time", v)
	}
	return time.Time(d).Equal(v2) || time.Time(d).Before(v2), nil
}

func (d TimestampValue) ToInt64() (int64, error) {
	return time.Time(d).Unix(), nil
}

func (d TimestampValue) ToString() (string, error) {
	json, err := d.ToJSON()
	if err != nil {
		return "", err
	}
	return toTimestampValueFromString(json)
}

func (d TimestampValue) ToBytes() ([]byte, error) {
	json, err := d.ToJSON()
	if err != nil {
		return nil, err
	}
	v, err := toTimestampValueFromString(json)
	if err != nil {
		return nil, err
	}
	return []byte(v), nil
}

func (d TimestampValue) ToFloat64() (float64, error) {
	return float64(time.Time(d).Unix()), nil
}

func (d TimestampValue) ToBool() (bool, error) {
	return false, fmt.Errorf("failed to convert %v to bool type", d)
}

func (d TimestampValue) ToArray() (*ArrayValue, error) {
	return nil, fmt.Errorf("failed to convert %v to array type", d)
}

func (d TimestampValue) ToStruct() (*StructValue, error) {
	return nil, fmt.Errorf("failed to convert %v to struct type", d)
}

func (d TimestampValue) ToJSON() (string, error) {
	return time.Time(d).Format(time.RFC3339Nano), nil
}

func (d TimestampValue) ToTime() (time.Time, error) {
	return time.Time(d), nil
}

func (d TimestampValue) ToRat() (*big.Rat, error) {
	return nil, fmt.Errorf("failed to convert *big.Rat from timestamp %v", d)
}

func (d TimestampValue) Marshal() (string, error) {
	json, err := d.ToJSON()
	if err != nil {
		return "", err
	}
	return toTimestampValueFromString(json)
}

func (d TimestampValue) Format(verb rune) string {
	formatted := time.Time(d).Format(time.RFC3339)
	switch verb {
	case 't':
		return formatted
	case 'T':
		return fmt.Sprintf(`TIMESTAMP "%s"`, formatted)
	}
	return formatted
}

func (d TimestampValue) Interface() interface{} {
	return time.Time(d).Format(time.RFC3339)
}

type SafeValue struct {
	value Value
}

func (v *SafeValue) Add(arg Value) (Value, error) {
	ret, err := v.value.Add(arg)
	if err != nil {
		return nil, nil
	}
	return ret, nil
}

func (v *SafeValue) Sub(arg Value) (Value, error) {
	ret, err := v.value.Sub(arg)
	if err != nil {
		return nil, nil
	}
	return ret, nil
}

func (v *SafeValue) Mul(arg Value) (Value, error) {
	ret, err := v.value.Mul(arg)
	if err != nil {
		return nil, nil
	}
	return ret, nil
}

func (v *SafeValue) Div(arg Value) (Value, error) {
	ret, err := v.value.Div(arg)
	if err != nil {
		return nil, nil
	}
	return ret, nil
}

func (v *SafeValue) EQ(arg Value) (bool, error) {
	ret, err := v.value.EQ(arg)
	if err != nil {
		return false, nil
	}
	return ret, nil
}

func (v *SafeValue) GT(arg Value) (bool, error) {
	ret, err := v.value.GT(arg)
	if err != nil {
		return false, nil
	}
	return ret, nil
}

func (v *SafeValue) GTE(arg Value) (bool, error) {
	ret, err := v.value.GTE(arg)
	if err != nil {
		return false, nil
	}
	return ret, nil
}

func (v *SafeValue) LT(arg Value) (bool, error) {
	ret, err := v.value.LT(arg)
	if err != nil {
		return false, nil
	}
	return ret, nil
}

func (v *SafeValue) LTE(arg Value) (bool, error) {
	ret, err := v.value.LTE(arg)
	if err != nil {
		return false, nil
	}
	return ret, nil
}

func (v *SafeValue) ToInt64() (int64, error) {
	ret, err := v.value.ToInt64()
	if err != nil {
		return 0, nil
	}
	return ret, nil
}

func (v *SafeValue) ToString() (string, error) {
	ret, err := v.value.ToString()
	if err != nil {
		return "", nil
	}
	return ret, nil
}

func (v *SafeValue) ToBytes() ([]byte, error) {
	ret, err := v.value.ToBytes()
	if err != nil {
		return nil, nil
	}
	return ret, nil
}

func (v *SafeValue) ToFloat64() (float64, error) {
	ret, err := v.value.ToFloat64()
	if err != nil {
		return 0, nil
	}
	return ret, nil
}

func (v *SafeValue) ToBool() (bool, error) {
	ret, err := v.value.ToBool()
	if err != nil {
		return false, nil
	}
	return ret, nil
}

func (v *SafeValue) ToArray() (*ArrayValue, error) {
	ret, err := v.value.ToArray()
	if err != nil {
		return &ArrayValue{}, nil
	}
	return ret, nil
}

func (v *SafeValue) ToStruct() (*StructValue, error) {
	ret, err := v.value.ToStruct()
	if err != nil {
		return &StructValue{}, nil
	}
	return ret, nil
}

func (v *SafeValue) ToJSON() (string, error) {
	ret, err := v.value.ToJSON()
	if err != nil {
		return "", nil
	}
	return ret, nil
}

func (v *SafeValue) ToTime() (time.Time, error) {
	ret, err := v.value.ToTime()
	if err != nil {
		return time.Time{}, nil
	}
	return ret, nil
}

func (v *SafeValue) ToRat() (*big.Rat, error) {
	ret, err := v.value.ToRat()
	if err != nil {
		return nil, nil
	}
	return ret, nil
}

func (v *SafeValue) Marshal() (string, error) {
	ret, err := v.value.Marshal()
	if err != nil {
		return "", nil
	}
	return ret, nil
}

func (v *SafeValue) Format(verb rune) string {
	return v.value.Format(verb)
}

func (v *SafeValue) Interface() interface{} {
	return v.value.Interface()
}

const (
	ArrayValueHeader     = "zetasqlitearray:"
	StructValueHeader    = "zetasqlitestruct:"
	BytesValueHeader     = "zetasqlitebytes:"
	NumericValueHeader   = "zetasqlitenumeric:"
	DateValueHeader      = "zetasqlitedate:"
	DatetimeValueHeader  = "zetasqlitedatetime:"
	TimeValueHeader      = "zetasqlitetime:"
	TimestampValueHeader = "zetasqlitetimestamp:"
	JsonValueHeader      = "zetasqlitejson:"
)

func ValueOf(v interface{}) (Value, error) {
	if v == nil {
		return nil, nil
	}
	if _, ok := v.([]byte); ok {
		if reflect.ValueOf(v).IsNil() {
			return nil, nil
		}
	}
	switch vv := v.(type) {
	case int:
		return IntValue(int64(vv)), nil
	case int8:
		return IntValue(int64(vv)), nil
	case int16:
		return IntValue(int64(vv)), nil
	case int32:
		return IntValue(int64(vv)), nil
	case int64:
		return IntValue(vv), nil
	case uint:
		return IntValue(int64(vv)), nil
	case uint8:
		return IntValue(int64(vv)), nil
	case uint16:
		return IntValue(int64(vv)), nil
	case uint32:
		return IntValue(int64(vv)), nil
	case uint64:
		return IntValue(int64(vv)), nil
	case string:
		switch {
		case isArrayValue(vv):
			return ArrayValueOf(vv)
		case isStructValue(vv):
			return StructValueOf(vv)
		case isBytesValue(vv):
			return BytesValueOf(vv)
		case isNumericValue(vv):
			return NumericValueOf(vv)
		case isDateValue(vv):
			return DateValueOf(vv)
		case isDatetimeValue(vv):
			return DatetimeValueOf(vv)
		case isTimeValue(vv):
			return TimeValueOf(vv)
		case isTimestampValue(vv):
			return TimestampValueOf(vv)
		case isJsonValue(vv):
			return JsonValueOf(vv)
		}
		return StringValue(vv), nil
	case []byte:
		return StringValue(string(vv)), nil
	case float32:
		return FloatValue(float64(vv)), nil
	case float64:
		return FloatValue(vv), nil
	case bool:
		return BoolValue(vv), nil
	}
	return nil, fmt.Errorf("failed to convert value from %T", v)
}

func isArrayValue(v string) bool {
	if len(v) < len(ArrayValueHeader) {
		return false
	}
	if v[0] == '"' {
		return strings.HasPrefix(v[1:], ArrayValueHeader)
	}
	return strings.HasPrefix(v, ArrayValueHeader)
}

func isStructValue(v string) bool {
	if len(v) < len(StructValueHeader) {
		return false
	}
	if v[0] == '"' {
		return strings.HasPrefix(v[1:], StructValueHeader)
	}
	return strings.HasPrefix(v, StructValueHeader)
}

func isBytesValue(v string) bool {
	if len(v) < len(BytesValueHeader) {
		return false
	}
	if v[0] == '"' {
		return strings.HasPrefix(v[1:], BytesValueHeader)
	}
	return strings.HasPrefix(v, BytesValueHeader)
}

func isNumericValue(v string) bool {
	if len(v) < len(NumericValueHeader) {
		return false
	}
	if v[0] == '"' {
		return strings.HasPrefix(v[1:], NumericValueHeader)
	}
	return strings.HasPrefix(v, NumericValueHeader)
}

func isDateValue(v string) bool {
	if len(v) < len(DateValueHeader) {
		return false
	}
	if v[0] == '"' {
		return strings.HasPrefix(v[1:], DateValueHeader)
	}
	return strings.HasPrefix(v, DateValueHeader)
}

func isDatetimeValue(v string) bool {
	if len(v) < len(DatetimeValueHeader) {
		return false
	}
	if v[0] == '"' {
		return strings.HasPrefix(v[1:], DatetimeValueHeader)
	}
	return strings.HasPrefix(v, DatetimeValueHeader)
}

func isTimeValue(v string) bool {
	if len(v) < len(TimeValueHeader) {
		return false
	}
	if v[0] == '"' {
		return strings.HasPrefix(v[1:], TimeValueHeader)
	}
	return strings.HasPrefix(v, TimeValueHeader)
}

func isTimestampValue(v string) bool {
	if len(v) < len(TimestampValueHeader) {
		return false
	}
	if v[0] == '"' {
		return strings.HasPrefix(v[1:], TimestampValueHeader)
	}
	return strings.HasPrefix(v, TimestampValueHeader)
}

func isJsonValue(v string) bool {
	if len(v) < len(JsonValueHeader) {
		return false
	}
	if v[0] == '"' {
		return strings.HasPrefix(v[1:], JsonValueHeader)
	}
	return strings.HasPrefix(v, JsonValueHeader)
}

func BytesValueOf(v string) (Value, error) {
	bytes, err := bytesValueFromEncodedString(v)
	if err != nil {
		return nil, fmt.Errorf("failed to get bytes value from encoded string: %w", err)
	}
	return BytesValue(bytes), nil
}

func NumericValueOf(v string) (Value, error) {
	numeric, err := numericValueFromEncodedString(v)
	if err != nil {
		return nil, fmt.Errorf("failed to get numeric value from encoded string: %w", err)
	}
	return (*NumericValue)(numeric), nil
}

func DateValueOf(v string) (Value, error) {
	date, err := dateValueFromEncodedString(v)
	if err != nil {
		return nil, fmt.Errorf("failed to get date value from encoded string: %w", err)
	}
	return DateValue(date), nil
}

func DatetimeValueOf(v string) (Value, error) {
	date, err := datetimeValueFromEncodedString(v)
	if err != nil {
		return nil, fmt.Errorf("failed to get datetime value from encoded string: %w", err)
	}
	return DatetimeValue(date), nil
}

func TimeValueOf(v string) (Value, error) {
	date, err := timeValueFromEncodedString(v)
	if err != nil {
		return nil, fmt.Errorf("failed to get time value from encoded string: %w", err)
	}
	return TimeValue(date), nil
}

func TimestampValueOf(v string) (Value, error) {
	date, err := timestampValueFromEncodedString(v)
	if err != nil {
		return nil, fmt.Errorf("failed to get timestamp value from encoded string: %w", err)
	}
	return TimestampValue(date), nil
}

func JsonValueOf(v string) (Value, error) {
	json, err := jsonValueFromEncodedString(v)
	if err != nil {
		return nil, fmt.Errorf("failed to get json value from encoded string: %w", err)
	}
	return JsonValue(json), nil
}

func ArrayValueOf(v string) (Value, error) {
	arr, err := arrayValueFromEncodedString(v)
	if err != nil {
		return nil, fmt.Errorf("failed to get array value from encoded string: %w", err)
	}
	values := make([]Value, 0, len(arr))
	for _, a := range arr {
		val, err := ValueOf(a)
		if err != nil {
			return nil, err
		}
		values = append(values, val)
	}
	return &ArrayValue{values: values}, nil
}

func StructValueOf(v string) (Value, error) {
	if len(v) == 0 {
		return nil, nil
	}
	if v[0] == '"' {
		unquoted, err := strconv.Unquote(v)
		if err != nil {
			return nil, fmt.Errorf("failed to unquote value %q: %w", v, err)
		}
		v = unquoted
	}
	content := v[len(StructValueHeader):]
	decoded, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return nil, fmt.Errorf("failed to base64 decode for struct value %q: %w", content, err)
	}
	dec := json.NewDecoder(bytes.NewBuffer(decoded))
	t, err := dec.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to decode struct value %q: %w", decoded, err)
	}
	if t != json.Delim('{') {
		return nil, fmt.Errorf("invalid delimiter of struct value %q", decoded)
	}
	var (
		keys   []string
		values []Value
		valMap = map[string]Value{}
	)
	for {
		k, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed to decode struct key %q: %w", decoded, err)
		}
		if k == json.Delim('}') {
			break
		}
		key := k.(string)
		var value interface{}
		if err := dec.Decode(&value); err != nil {
			return nil, fmt.Errorf("failed to decode struct value %q: %w", decoded, err)
		}
		keys = append(keys, key)
		val, err := ValueOf(value)
		if err != nil {
			return nil, fmt.Errorf("failed to convert value from %v: %w", value, err)
		}
		values = append(values, val)
		valMap[key] = val
	}
	return &StructValue{keys: keys, values: values, m: valMap}, nil
}

func toArrayValueFromJSONString(json string) string {
	return strconv.Quote(
		fmt.Sprintf(
			"%s%s",
			ArrayValueHeader,
			base64.StdEncoding.EncodeToString([]byte(json)),
		),
	)
}

func bytesValueFromEncodedString(v string) ([]byte, error) {
	if len(v) == 0 {
		return nil, nil
	}
	if v[0] == '"' {
		unquoted, err := strconv.Unquote(v)
		if err != nil {
			return nil, fmt.Errorf("failed to unquote value %q: %w", v, err)
		}
		v = unquoted
	}
	content := v[len(BytesValueHeader):]
	decoded, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

func numericValueFromEncodedString(v string) (*big.Rat, error) {
	if len(v) == 0 {
		return nil, nil
	}
	if v[0] == '"' {
		unquoted, err := strconv.Unquote(v)
		if err != nil {
			return nil, fmt.Errorf("failed to unquote value %q: %w", v, err)
		}
		v = unquoted
	}
	content := v[len(NumericValueHeader):]
	ret := new(big.Rat)
	ret.SetString(content)
	return ret, nil
}

func dateValueFromEncodedString(v string) (time.Time, error) {
	if len(v) == 0 {
		return time.Time{}, nil
	}
	if v[0] == '"' {
		unquoted, err := strconv.Unquote(v)
		if err != nil {
			return time.Time{}, fmt.Errorf("failed to unquote value %q: %w", v, err)
		}
		v = unquoted
	}
	content := v[len(DateValueHeader):]
	return parseDate(content)
}

func datetimeValueFromEncodedString(v string) (time.Time, error) {
	if len(v) == 0 {
		return time.Time{}, nil
	}
	if v[0] == '"' {
		unquoted, err := strconv.Unquote(v)
		if err != nil {
			return time.Time{}, fmt.Errorf("failed to unquote value %q: %w", v, err)
		}
		v = unquoted
	}
	content := v[len(DatetimeValueHeader):]
	return parseDatetime(content)
}

func timeValueFromEncodedString(v string) (time.Time, error) {
	if len(v) == 0 {
		return time.Time{}, nil
	}
	if v[0] == '"' {
		unquoted, err := strconv.Unquote(v)
		if err != nil {
			return time.Time{}, fmt.Errorf("failed to unquote value %q: %w", v, err)
		}
		v = unquoted
	}
	content := v[len(TimeValueHeader):]
	return parseTime(content)
}

func timestampValueFromEncodedString(v string) (time.Time, error) {
	if len(v) == 0 {
		return time.Time{}, nil
	}
	if v[0] == '"' {
		unquoted, err := strconv.Unquote(v)
		if err != nil {
			return time.Time{}, fmt.Errorf("failed to unquote value %q: %w", v, err)
		}
		v = unquoted
	}
	content := v[len(TimestampValueHeader):]
	loc, err := time.LoadLocation("")
	if err != nil {
		return time.Time{}, err
	}
	return parseTimestamp(content, loc)
}

func jsonValueFromEncodedString(v string) (string, error) {
	if len(v) == 0 {
		return "", nil
	}
	if v[0] == '"' {
		unquoted, err := strconv.Unquote(v)
		if err != nil {
			return "", fmt.Errorf("failed to unquote value %q: %w", v, err)
		}
		v = unquoted
	}
	content := v[len(JsonValueHeader):]
	decoded, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return "", fmt.Errorf("failed to base64 decode for json value %q: %w", content, err)
	}
	return string(decoded), nil
}

func arrayValueFromEncodedString(v string) ([]interface{}, error) {
	if len(v) == 0 {
		return nil, nil
	}
	if v[0] == '"' {
		unquoted, err := strconv.Unquote(v)
		if err != nil {
			return nil, fmt.Errorf("failed to unquote value %q: %w", v, err)
		}
		v = unquoted
	}
	content := v[len(ArrayValueHeader):]
	decoded, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return nil, fmt.Errorf("failed to base64 decode for array value %q: %w", content, err)
	}
	var arr []interface{}
	if err := json.Unmarshal(decoded, &arr); err != nil {
		return nil, fmt.Errorf("failed to decode array: %w", err)
	}
	return arr, nil
}

func jsonArrayFromEncodedString(v string) ([]byte, error) {
	if len(v) == 0 {
		return nil, nil
	}
	if v[0] == '"' {
		unquoted, err := strconv.Unquote(v)
		if err != nil {
			return nil, fmt.Errorf("failed to unquote value %q: %w", v, err)
		}
		v = unquoted
	}
	content := v[len(ArrayValueHeader):]
	decoded, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return nil, fmt.Errorf("failed to base64 decode for array value %q: %w", content, err)
	}
	return decoded, nil
}

func toStructValueFromJSONString(json string) string {
	return strconv.Quote(
		fmt.Sprintf(
			"%s%s",
			StructValueHeader,
			base64.StdEncoding.EncodeToString([]byte(json)),
		),
	)
}

func toBytesValueFromString(s string) string {
	return strconv.Quote(
		fmt.Sprintf(
			"%s%s",
			BytesValueHeader,
			base64.StdEncoding.EncodeToString([]byte(s)),
		),
	)
}

func toNumericValueFromString(s string) string {
	return strconv.Quote(
		fmt.Sprintf(
			"%s%s",
			NumericValueHeader,
			s,
		),
	)
}

func toDateValueFromString(s string) string {
	return strconv.Quote(
		fmt.Sprintf(
			"%s%s",
			DateValueHeader,
			s,
		),
	)
}

func toDatetimeValueFromString(s string) string {
	return strconv.Quote(
		fmt.Sprintf(
			"%s%s",
			DatetimeValueHeader,
			s,
		),
	)
}

func toTimeValueFromString(s string) string {
	return strconv.Quote(
		fmt.Sprintf(
			"%s%s",
			TimeValueHeader,
			s,
		),
	)
}

func toTimestampValueFromString(s string) (string, error) {
	formatted, err := formatTimestamp(s)
	if err != nil {
		return "", err
	}
	return strconv.Quote(
		fmt.Sprintf(
			"%s%s",
			TimestampValueHeader,
			formatted,
		),
	), nil
}

func toJsonValueFromString(s string) (string, error) {
	return strconv.Quote(
		fmt.Sprintf(
			"%s%s",
			JsonValueHeader,
			base64.StdEncoding.EncodeToString([]byte(s)),
		),
	), nil
}

func formatTimestamp(s string) (string, error) {
	loc, err := time.LoadLocation("")
	if err != nil {
		return "", err
	}
	t, err := parseTimestamp(s, loc)
	if err != nil {
		return "", err
	}
	return t.Format(time.RFC3339Nano), nil
}

func isNULLValue(v interface{}) bool {
	vv, ok := v.([]byte)
	if !ok {
		return false
	}
	return vv == nil
}

var (
	dateRe     = regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2}$`)
	datetimeRe = regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2}[T\s][0-9]{2}:[0-9]{2}:[0-9]{2}$`)
	timeRe     = regexp.MustCompile(`^[0-9]{2}:[0-9]{2}:[0-9]{2}$`)
)

func isDate(date string) bool {
	return dateRe.MatchString(date)
}

func isDatetime(datetime string) bool {
	return datetimeRe.MatchString(datetime)
}

func isTime(time string) bool {
	return timeRe.MatchString(time)
}

func isTimestamp(timestamp string) bool {
	loc, err := time.LoadLocation("")
	if err != nil {
		return false
	}
	if _, err := parseTimestamp(timestamp, loc); err != nil {
		return false
	}
	return true
}

func parseDate(date string) (time.Time, error) {
	return time.Parse("2006-01-02", date)
}

func parseDatetime(datetime string) (time.Time, error) {
	if t, err := time.Parse("2006-01-02T15:04:05", datetime); err == nil {
		return t, nil
	}
	return time.Parse("2006-01-02 15:04:05", datetime)
}

func parseTime(t string) (time.Time, error) {
	return time.Parse("15:04:05", t)
}

func parseTimestamp(timestamp string, loc *time.Location) (time.Time, error) {
	if t, err := time.ParseInLocation("2006-01-02T15:04:05.999999999Z07:00", timestamp, loc); err == nil {
		return t, nil
	}
	if t, err := time.ParseInLocation("2006-01-02T15:04:05.999999999-07:00", timestamp, loc); err == nil {
		return t, nil
	}
	if t, err := time.ParseInLocation("2006-01-02T15:04:05.999999999-07", timestamp, loc); err == nil {
		return t, nil
	}
	if t, err := time.ParseInLocation("2006-01-02T15:04:05.999999999 MST", timestamp, loc); err == nil {
		return t, nil
	}
	if t, err := time.ParseInLocation("2006-01-02T15:04:05Z07:00", timestamp, loc); err == nil {
		return t, nil
	}
	if t, err := time.ParseInLocation("2006-01-02T15:04:05-07:00", timestamp, loc); err == nil {
		return t, nil
	}
	if t, err := time.ParseInLocation("2006-01-02T15:04:05-07", timestamp, loc); err == nil {
		return t, nil
	}
	if t, err := time.ParseInLocation("2006-01-02T15:04:05 MST", timestamp, loc); err == nil {
		return t, nil
	}
	if t, err := time.ParseInLocation("2006-01-02T15:04:05", timestamp, loc); err == nil {
		return t, nil
	}
	if t, err := time.ParseInLocation("2006-01-02 15:04:05.999999999Z07:00", timestamp, loc); err == nil {
		return t, nil
	}
	if t, err := time.ParseInLocation("2006-01-02 15:04:05.999999999-07:00", timestamp, loc); err == nil {
		return t, nil
	}
	if t, err := time.ParseInLocation("2006-01-02 15:04:05.999999999-07", timestamp, loc); err == nil {
		return t, nil
	}
	if t, err := time.ParseInLocation("2006-01-02 15:04:05.999999999 MST", timestamp, loc); err == nil {
		return t, nil
	}
	if t, err := time.ParseInLocation("2006-01-02 15:04:05Z07:00", timestamp, loc); err == nil {
		return t, nil
	}
	if t, err := time.ParseInLocation("2006-01-02 15:04:05+07:00", timestamp, loc); err == nil {
		return t, nil
	}
	if t, err := time.ParseInLocation("2006-01-02 15:04:05-07", timestamp, loc); err == nil {
		return t, nil
	}
	if t, err := time.ParseInLocation("2006-01-02 15:04:05 MST", timestamp, loc); err == nil {
		return t, nil
	}
	if t, err := time.ParseInLocation("2006-01-02 15:04:05", timestamp, loc); err == nil {
		return t, nil
	}
	if t, err := time.ParseInLocation("2006-01-02", timestamp, loc); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("failed to parse timestamp. unexpected format %s", timestamp)
}

func toDateValueFromInt64(days int64) string {
	t := time.Unix(int64(time.Duration(days)*24*time.Hour/time.Second), 0)
	return t.Format("2006-01-02")
}

const (
	microSecondShift = 20
	secShift         = 0
	minShift         = 6
	hourShift        = 12
	dayShift         = 17
	monthShift       = 22
	yearShift        = 26
	secMask          = 0b111111
	minMask          = 0b111111 << minShift
	hourMask         = 0b11111 << hourShift
	dayMask          = 0b11111 << dayShift
	monthMask        = 0b1111 << monthShift
	yearMask         = 0x3FFF << yearShift
)

func toDatetimeValueFromInt64(bit int64) string {
	b := bit >> 20
	year := (b & yearMask) >> yearShift
	month := (b & monthMask) >> monthShift
	day := (b & dayMask) >> dayShift
	hour := (b & hourMask) >> hourShift
	min := (b & minMask) >> minShift
	sec := (b & secMask) >> secShift
	return fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d", year, month, day, hour, min, sec)
}

func toTimeValueFromInt64(bit int64) string {
	b := bit >> 20
	hour := (b & hourMask) >> hourShift
	min := (b & minMask) >> minShift
	sec := (b & secMask) >> secShift
	return fmt.Sprintf("%02d:%02d:%02d", hour, min, sec)
}

func toTimestampValueFromTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

func toTimeValue(s string) (time.Time, error) {
	switch {
	case isTimestamp(s):
		loc, err := time.LoadLocation("")
		if err != nil {
			return time.Time{}, err
		}
		return parseTimestamp(s, loc)
	case isDatetime(s):
		return parseDatetime(s)
	case isDate(s):
		return parseDate(s)
	case isTime(s):
		return parseTime(s)
	}
	return time.Time{}, fmt.Errorf("unsupported time format %s", s)
}

func EncodeNamedValues(v []driver.NamedValue, params []*ast.ParameterNode) ([]sql.NamedArg, error) {
	if len(v) != len(params) {
		return nil, fmt.Errorf(
			"failed to match named values num (%d) and params num (%d)",
			len(v), len(params),
		)
	}
	ret := make([]sql.NamedArg, 0, len(v))
	for idx, vv := range v {
		converted, err := encodeNamedValue(vv, params[idx])
		if err != nil {
			return nil, fmt.Errorf("failed to convert value from %+v: %w", vv, err)
		}
		ret = append(ret, converted)
	}
	return ret, nil
}

func encodeNamedValue(v driver.NamedValue, param *ast.ParameterNode) (sql.NamedArg, error) {
	value, err := encodeValueWithType(v.Value, param.Type())
	if err != nil {
		return sql.NamedArg{}, err
	}
	return sql.NamedArg{
		Name:  strings.ToLower(v.Name),
		Value: value,
	}, nil
}

func encodeValues(v []interface{}, params []*ast.ParameterNode) ([]interface{}, error) {
	if len(v) != len(params) {
		return nil, fmt.Errorf(
			"failed to match args values num (%d) and params num (%d)",
			len(v), len(params),
		)
	}
	ret := make([]interface{}, 0, len(v))
	for idx, vv := range v {
		value, err := encodeValueWithType(vv, params[idx].Type())
		if err != nil {
			return nil, err
		}
		ret = append(ret, value)
	}
	return ret, nil
}

func encodeValueWithType(v interface{}, t types.Type) (interface{}, error) {
	if v == nil {
		return nil, nil
	}
	switch t.Kind() {
	case types.INT32, types.INT64, types.UINT32, types.UINT64, types.ENUM:
		vv, err := ValueOf(v)
		if err != nil {
			return nil, err
		}
		return vv.ToInt64()
	case types.BOOL:
		vv, err := ValueOf(v)
		if err != nil {
			return nil, err
		}
		return vv.ToBool()
	case types.FLOAT, types.DOUBLE:
		vv, err := ValueOf(v)
		if err != nil {
			return nil, err
		}
		return vv.ToFloat64()
	case types.STRING:
		vv, err := ValueOf(v)
		if err != nil {
			return nil, err
		}
		return vv.ToString()
	case types.BYTES:
		vv, err := ValueOf(v)
		if err != nil {
			return nil, err
		}
		return vv.ToString()
	case types.DATE:
		text, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("failed to convert DATE from %T", v)
		}
		return toDateValueFromString(text), nil
	case types.TIMESTAMP:
		text, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("failed to convert TIMESTAMP from %T", v)
		}
		return toTimestampValueFromString(text)
	case types.ARRAY:
		rv := reflect.ValueOf(v)
		if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
			return nil, fmt.Errorf("failed to convert %v to array", v)
		}
		elemType := t.AsArray().ElementType()
		var values []string
		for i := 0; i < rv.Len(); i++ {
			value, err := encodeValueWithType(rv.Index(i).Interface(), elemType)
			if err != nil {
				return nil, err
			}
			if value == nil {
				values = append(values, "null")
			} else if elemType.Kind() == types.STRING || elemType.Kind() == types.BYTES {
				values = append(values, strconv.Quote(fmt.Sprint(value)))
			} else {
				values = append(values, fmt.Sprint(value))
			}
		}
		value := fmt.Sprintf("[%s]", strings.Join(values, ","))
		return toArrayValueFromJSONString(value), nil
	case types.STRUCT:
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Ptr:
			return encodeValueWithType(rv.Elem().Interface(), t)
		case reflect.Struct:
			st := t.AsStruct()
			structType := rv.Type()
			fields := make([]string, 0, rv.NumField())
			for i := 0; i < rv.NumField(); i++ {
				typ := st.Field(i).Type()
				field, err := encodeValueWithType(rv.Field(i), typ)
				if err != nil {
					return nil, err
				}
				if field == nil {
					fields = append(fields, fmt.Sprintf(`"%s":null`))
				} else if typ.Kind() == types.STRING || typ.Kind() == types.BYTES {
					fields = append(fields, strconv.Quote(fmt.Sprint(field)))
				} else {
					fields = append(fields, fmt.Sprintf(`"%s":%s`, structType.Field(i), fmt.Sprint(field)))
				}
			}
			value := fmt.Sprintf("{%s}", strings.Join(fields, ","))
			return toStructValueFromJSONString(value), nil
		case reflect.Slice:
			// we expect []map[string]interface{} type for struct
			st := t.AsStruct()
			if st.NumFields() != rv.Len() {
				return nil, fmt.Errorf("unexpected field number. expected %d but got %d", st.NumFields(), rv.Len())
			}
			fields := make([]string, 0, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				elem := rv.Index(i)
				if elem.Kind() != reflect.Map {
					return nil, fmt.Errorf("unexpected element type of slice %s for struct type. please use map[string]interface{}", elem.Kind())
				}
				keys := elem.MapKeys()
				if len(keys) != 1 {
					return nil, fmt.Errorf("unexpected map key number. expected one key for column but got %d", len(keys))
				}
				if keys[0].Kind() != reflect.String {
					return nil, fmt.Errorf("unexpected map key type. expected string type but got %s", keys[0].Kind())
				}
				fieldName := keys[0].Interface().(string)
				fieldValue := elem.MapIndex(keys[0]).Interface()
				field, err := encodeValueWithType(fieldValue, st.Field(i).Type())
				if err != nil {
					return nil, err
				}
				fields = append(fields, fmt.Sprintf(`"%s":%s`, fieldName, fmt.Sprint(field)))
			}
			value := fmt.Sprintf("{%s}", strings.Join(fields, ","))
			return toStructValueFromJSONString(value), nil
		case reflect.Map:
			return nil, fmt.Errorf("unsupported map type for STRUCT column. please use slice or struct type")
		default:
			return nil, fmt.Errorf("failed to convert %v to struct", v)
		}
	case types.TIME:
		text, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("failed to convert TIME from %T", v)
		}
		return toTimeValueFromString(text), nil
	case types.DATETIME:
		text, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("failed to convert DATETIME from %T", v)
		}
		return toDatetimeValueFromString(text), nil
	case types.PROTO:
		return nil, fmt.Errorf("failed to convert PROTO type from %T", v)
	case types.GEOGRAPHY:
		return nil, fmt.Errorf("failed to convert GEOGRAPHY type from %T", v)
	case types.NUMERIC:
		return nil, fmt.Errorf("failed to convert NUMERIC type from %T", v)
	case types.BIG_NUMERIC:
		return nil, fmt.Errorf("failed to convert BIGNUMERIC type from %T", v)
	case types.EXTENDED:
		return nil, fmt.Errorf("failed to convert EXTENDED type from %T", v)
	case types.JSON:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		return string(b), nil
	case types.INTERVAL:
		return nil, fmt.Errorf("failed to convert INTERVAL type from %T", v)
	default:
	}
	return nil, fmt.Errorf("unexpected type %s to convert from %T", t.Kind(), v)
}
