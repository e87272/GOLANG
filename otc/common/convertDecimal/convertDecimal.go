package convertDecimal

import (
	"math"
	"regexp"
	"strconv"
	"strings"
)

// 十進位小數字串與按照 DecimalDigit 與整數互相轉換，避免浮點誤差

const DecimalDigit = 6

type Decimal struct {
	value int64
}

var Zero = Decimal{value: int64(0)}

// 乘上多少倍可以轉換成整數
var multipleFloat = math.Pow10(DecimalDigit)
var multipleInt = int64(multipleFloat)

func FromInt64(valueInt64 int64) Decimal {
	return Decimal{value: valueInt64 * multipleInt}
}

func FromFloat64(valueFloat64 float64) Decimal {
	return Decimal{value: int64(math.Round(valueFloat64 * multipleFloat))}
}

func FromString(valueStr string) (Decimal, bool) {

	// 十進位格式判斷
	ok, _ := regexp.MatchString(`^-?[0-9]+(\.[0-9]*)?$`, valueStr)
	if !ok {
		// 不符合十進位格式
		return Zero, false
	}

	// 補0到指定的小數位數，並去掉小數點
	a := strings.Split(valueStr, ".")
	allDigitsStr := a[0]
	if len(a) > 1 {
		fractionStr := strings.TrimRight(a[1], "0")
		if len(fractionStr) > DecimalDigit {
			// 小數位數超過上限
			return Zero, false
		}
		allDigitsStr += fractionStr + strings.Repeat("0", DecimalDigit-len(fractionStr))
	} else {
		allDigitsStr += strings.Repeat("0", DecimalDigit)
	}

	// 轉換成int64
	value, err := strconv.ParseInt(allDigitsStr, 10, 64)
	if err != nil {
		// 超過int64上限
		return Zero, false
	}

	return Decimal{value: value}, true
}

func (decimal Decimal) ToInt64() int64 {
	return decimal.value / multipleInt
}

func (decimal Decimal) ToFloat64() float64 {
	return float64(decimal.value) / multipleFloat
}

func (decimal Decimal) ToString() string {

	// 轉成字串，分成數字部分和正負號
	allDigitsStr := strconv.FormatInt(decimal.value, 10)
	valueStr := ""
	if allDigitsStr[:1] == "-" {
		valueStr = "-"
		allDigitsStr = allDigitsStr[1:]
	}

	// 前面補0到超過指定小數位數
	if len(allDigitsStr) <= DecimalDigit {
		allDigitsStr = strings.Repeat("0", 1+DecimalDigit-len(allDigitsStr)) + allDigitsStr
	}

	// 指定位數以上的為整數位，以下的為小數位
	valueStr += allDigitsStr[:len(allDigitsStr)-DecimalDigit]
	fractionStr := strings.TrimRight(allDigitsStr[len(allDigitsStr)-DecimalDigit:], "0")
	if len(fractionStr) > 0 {
		valueStr += "." + fractionStr
	}

	return valueStr
}

func (decimal1 Decimal) GreaterThan(decimal2 Decimal) bool {
	return decimal1.value > decimal2.value
}

func (decimal1 Decimal) LessThan(decimal2 Decimal) bool {
	return decimal1.value < decimal2.value
}

func (decimal1 Decimal) Add(decimal2 Decimal) Decimal {
	return Decimal{value: decimal1.value + decimal2.value}
}

func (decimal1 Decimal) Sub(decimal2 Decimal) Decimal {
	return Decimal{value: decimal1.value - decimal2.value}
}

func (decimal1 Decimal) Mul(decimal2 Decimal) Decimal {

	// 分成整數部分和小數部分
	quo1 := decimal1.value / multipleInt
	rem1 := decimal1.value % multipleInt

	// 分成整數部分和小數部分
	quo2 := decimal2.value / multipleInt
	rem2 := decimal2.value % multipleInt

	return Decimal{value: decimal1.value*quo2 + quo1*rem2 + rem1*rem2/multipleInt}
}

func (decimal1 Decimal) Div(decimal2 Decimal) Decimal {

	// 分成整數部分和小數部分
	quo1 := decimal1.value / decimal2.value
	rem1 := decimal1.value % decimal2.value

	// 分成整數部分和小數部分
	quo2 := multipleInt / decimal2.value
	rem2 := multipleInt % decimal2.value

	return Decimal{value: decimal1.value*quo2 + quo1*rem2 + rem1*rem2/decimal2.value}
}

func (decimal1 Decimal) Mod(decimal2 Decimal) Decimal {
	return Decimal{value: decimal1.value % decimal2.value}
}
