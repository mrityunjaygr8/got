package got

import (
	"bytes"
	"errors"

	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type kvlm = orderedmap.OrderedMap[string, [][]byte]

func kvlm_parse(raw []byte, start int, dict kvlm) (kvlm, error) {
	var dct kvlm
	if dict.Len() == 0 {
		dct = *orderedmap.New[string, [][]byte]()
	} else {
		dct = dict
	}

	spc := bytes.IndexByte(raw[start:], ' ')
	nl := bytes.IndexByte(raw[start:], '\n')

	// base case
	if (spc < 0) || (nl < spc) {
		if !(nl == start) {
			return kvlm{}, nil
		}
		dct.Set("", [][]byte{raw[start+1:]})
		return dct, nil
	}

	key := string(raw[start:spc])
	end := start

	for {
		end = bytes.IndexByte(raw[end+1:], '\n')
		if raw[end+1] != ' ' {
			break
		}
	}

	value := bytes.ReplaceAll(raw[spc+1:end], []byte("\n "), []byte("\n"))
	old, present := dct.Get(key)
	if !present {
		dct.Set(key, [][]byte{value})
	} else {
		dct.Set(key, append(old, value))
	}

	return dct, nil
}

func kvln_serialize(kvlm kvlm) ([]byte, error) {
	raw := make([]byte, 0)

	for pair := kvlm.Oldest(); pair != nil; pair = pair.Next() {
		if pair.Key == string("") {
			continue
		}

		val := pair.Value
		for _, d := range val {
			data := make([]byte, 0)
			data = append(data, []byte(pair.Key)...)
			data = append(data, []byte(" ")...)
			data = append(data, bytes.ReplaceAll(d, []byte("\n"), []byte("\n "))...)
			data = append(data, []byte("\n")...)
			raw = append(raw, data...)
		}

	}

	raw = append(raw, []byte("\n")...)
	msg, present := kvlm.Get("")
	if !present {
		return nil, errors.New("There is no message")
	}
	raw = append(raw, msg[0]...)

	return raw, nil

}
