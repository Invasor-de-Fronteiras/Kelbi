package stringsupport

import (
	"bytes"
	"io/ioutil"
	"strconv"
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// StringConverter is a small helper for encoding/decoding strings.
type StringConverter struct {
	Encoding encoding.Encoding
}

// Decode decodes the given bytes as the set encoding.
func (sc *StringConverter) Decode(data []byte) (string, error) {
	decoded, err := ioutil.ReadAll(transform.NewReader(bytes.NewBuffer(data), sc.Encoding.NewDecoder()))

	if err != nil {
		return "", err
	}

	return string(decoded), nil
}

// MustDecode decodes the given bytes as the set encoding. Panics on decode failure.
func (sc *StringConverter) MustDecode(data []byte) string {
	decoded, err := sc.Decode(data)
	if err != nil {
		panic(err)
	}

	return decoded
}

// Encode encodes the given string as the set encoding.
func (sc *StringConverter) Encode(data string) ([]byte, error) {
	encoded, err := ioutil.ReadAll(transform.NewReader(bytes.NewBuffer([]byte(data)), sc.Encoding.NewEncoder()))

	if err != nil {
		return nil, err
	}

	return encoded, nil
}

// MustEncode encodes the given string as the set encoding. Panics on encode failure.
func (sc *StringConverter) MustEncode(data string) []byte {
	encoded, err := sc.Encode(data)
	if err != nil {
		panic(err)
	}

	return encoded
}

/*
func MustConvertShiftJISToUTF8(text string) string {
	result, err := ConvertShiftJISToUTF8(text)
	if err != nil {
		panic(err)
	}
	return result
}
func MustConvertUTF8ToShiftJIS(text string) string {
	result, err := ConvertUTF8ToShiftJIS(text)
	if err != nil {
		panic(err)
	}
	return result
}
func ConvertShiftJISToUTF8(text string) (string, error) {
	r := bytes.NewBuffer([]byte(text))
	decoded, err := ioutil.ReadAll(transform.NewReader(r, japanese.ShiftJIS.NewDecoder()))
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
*/

func UTF8ToSJIS(x string) []byte {
	e := japanese.ShiftJIS.NewEncoder()
	xt, _, err := transform.String(e, x)
	if err != nil {
		panic(err)
	}
	return []byte(xt)
}

func SJISToUTF8(b []byte) string {
	d := japanese.ShiftJIS.NewDecoder()
	result, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader(b), d))
	if err != nil {
		panic(err)
	}
	return string(result)
}

func PaddedString(x string, size uint, t bool) []byte {
	if t {
		e := japanese.ShiftJIS.NewEncoder()
		xt, _, err := transform.String(e, x)
		if err != nil {
			return make([]byte, size)
		}
		x = xt
	}
	out := make([]byte, size)
	copy(out, x)
	out[len(out)-1] = 0
	return out
}

func CSVAdd(csv string, v int) string {
	if len(csv) == 0 {
		return strconv.Itoa(v)
	}
	return csv + "," + strconv.Itoa(v)
}

func CSVRemove(csv string, v int) string {
	s := strings.Split(csv, ",")
	for i, e := range s {
		if e == strconv.Itoa(v) {
			s[i] = s[len(s)-1]
			s = s[:len(s)-1]
		}
	}
	return strings.Join(s, ",")
}

func CSVContains(csv string, v int) bool {
	s := strings.Split(csv, ",")
	for i := 0; i < len(s); i++ {
		j, _ := strconv.ParseInt(s[i], 10, 64)
		if int(j) == v {
			return true
		}
	}
	return false
}

func CSVLength(csv string) int {
	if csv == "" {
		return 0
	}
	s := strings.Split(csv, ",")
	return len(s)
}

func CSVElems(csv string) []int {
	var r []int
	if csv == "" {
		return r
	}
	s := strings.Split(csv, ",")
	for i := 0; i < len(s); i++ {
		j, _ := strconv.ParseInt(s[i], 10, 64)
		r = append(r, int(j))
	}
	return r
}

// ConvertUTF8ToShiftJIS converts a UTF8 string to a Shift-JIS []byte.
func ConvertUTF8ToShiftJIS(text string) ([]byte, error) {
	r := bytes.NewBuffer([]byte(text))
	encoded, err := ioutil.ReadAll(transform.NewReader(r, japanese.ShiftJIS.NewEncoder()))
	if err != nil {
		return nil, err
	}

	return encoded, nil
}

func ConvertUTF8ToSJIS(text string) (string, error) {
	r, _, err := transform.String(japanese.ShiftJIS.NewEncoder(), text)
	if err != nil {
		return "", err
	}
	return r, nil
}
