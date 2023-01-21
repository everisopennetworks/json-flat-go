package flat

import (
	"errors"
	"github.com/everisopennetworks/mergo"
	"reflect"
	"strconv"
	"strings"
)

// Options the flatten options.
// By default: Demiliter = "."
type Options struct {
	Delimiter          string
	Safe               bool
	MaxDepth           int
	ArrayDelimiter     string // PARENS (), BRACKETS [], CURLY_BRACES {}, "" NONE
	SliceDeepMerge     bool
	OverwriteNilInMaps bool
}

func getArrayDelimiters(str string) ([2]string, error) {
	switch str {
	case "PARENS":
		return [2]string{"(", ")"}, nil
	case "BRACKETS":
		return [2]string{"[", "]"}, nil
	case "CURLY_BRACES":
		return [2]string{"{", "}"}, nil
	default:
		return [2]string{}, errors.New("array delimiter not supported")
	}
}

// Flatten the map, it returns a map one level deep
// regardless of how nested the original map was.
// By default, the flatten has Delimiter = ".", and
// no limitation of MaxDepth
func Flatten(nested map[string]interface{}, opts *Options) (m map[string]interface{}, err error) {
	if opts == nil {
		opts = &Options{
			Delimiter: ".",
		}
	}

	m, err = flatten("", 0, nested, opts)

	return
}

func flatten(prefix string, depth int, nested interface{}, opts *Options) (flatmap map[string]interface{}, err error) {
	flatmap = make(map[string]interface{})

	switch nested := nested.(type) {
	case map[string]interface{}:
		if opts.MaxDepth != 0 && depth >= opts.MaxDepth {
			flatmap[prefix] = nested
			return
		}
		if reflect.DeepEqual(nested, map[string]interface{}{}) {
			flatmap[prefix] = nested
			return
		}
		for k, v := range nested {
			// create new key
			newKey := k
			if prefix != "" {
				newKey = prefix + opts.Delimiter + newKey
			}
			fm1, fe := flatten(newKey, depth+1, v, opts)
			if fe != nil {
				err = fe
				return
			}
			update(flatmap, fm1)
		}
	case []interface{}:
		if opts.Safe {
			flatmap[prefix] = nested
			return
		}
		if reflect.DeepEqual(nested, []interface{}{}) {
			flatmap[prefix] = nested
			return
		}
		for i, v := range nested {
			newKey := strconv.Itoa(i)

			ad, e := getArrayDelimiters(opts.ArrayDelimiter)
			if e != nil {
				ad = [2]string{"", ""}
			}
			newKey = ad[0] + newKey + ad[1]

			if prefix != "" && opts.ArrayDelimiter != "" {
				newKey = prefix + newKey
			} else if prefix != "" {
				newKey = prefix + opts.Delimiter + newKey
			}
			fm1, fe := flatten(newKey, depth+1, v, opts)
			if fe != nil {
				err = fe
				return
			}
			update(flatmap, fm1)
		}
	default:
		flatmap[prefix] = nested
	}
	return
}

// update is the function that update to map with from
// example:
// to = {"hi": "there"}
// from = {"foo": "bar"}
// then, to = {"hi": "there", "foo": "bar"}
func update(to map[string]interface{}, from map[string]interface{}) {
	for kt, vt := range from {
		to[kt] = vt
	}
}

// Unflatten the map, it returns a nested map of a map
// By default, the flatten has Delimiter = "."
func Unflatten(flat map[string]interface{}, opts *Options) (nested map[string]interface{}, err error) {
	if opts == nil {
		opts = &Options{
			Delimiter: ".",
		}
	}
	nested, err = unflatten(flat, opts)
	return
}

func unflatten(flat map[string]interface{}, opts *Options) (nested map[string]interface{}, err error) {
	nested = make(map[string]interface{})

	c := mergo.Config{}
	c.Overwrite = true
	c.OverwriteNilInMaps = opts.OverwriteNilInMaps

	if opts.SliceDeepMerge {
		mergo.WithSliceDeepMerge(&c)
		c.Overwrite = false
	}

	for k, v := range flat {
		temp := uf(k, v, opts).(map[string]interface{})
		err = mergo.Merge(&nested, temp, func(c2 *mergo.Config) {
			*c2 = c
		})
		if err != nil {
			return
		}
	}

	return
}

func getArrayIndex(key string, delim [2]string) (int, error) {
	if strings.HasPrefix(key, delim[0]) {
		val, err := strconv.Atoi(strings.Replace(strings.Replace(key, delim[0], "", -1), delim[1], "", -1))
		if err == nil {
			return val, nil
		} else {
			return -1, errors.New("key doesn't have int value")
		}
	} else {
		return -1, errors.New("key doesn't start with delimiter")
	}
}

func uf(k string, v interface{}, opts *Options) (n interface{}) {
	n = v

	ad, e := getArrayDelimiters(opts.ArrayDelimiter)
	if e == nil {
		k = strings.Replace(k, ad[0], opts.Delimiter+ad[0], -1)
	}

	keys := strings.Split(k, opts.Delimiter)

	for i := len(keys) - 1; i >= 0; i-- {
		idx, errIdx := getArrayIndex(keys[i], ad)
		if e != nil || (e == nil && errIdx != nil) {
			temp := make(map[string]interface{})
			temp[keys[i]] = n
			n = temp
		} else {
			temp := make([]interface{}, idx+1)
			temp[idx] = n
			n = temp
		}
	}

	return
}
