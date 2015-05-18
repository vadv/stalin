package ql

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"stalin/plugins/store/query_language/parser"
)

func Walk(node interface{}, data *Problems) (*Problems, error) {
	return walkNode(node, data)
}

func walkNode(node interface{}, data *Problems) (*Problems, error) {
	switch reflect.TypeOf(node) {
	case reflect.TypeOf(new(parser.NodeEqual)):
		return walkEqual(node.(*parser.NodeEqual), data)
	case reflect.TypeOf(new(parser.NodeLike)):
		return walkLike(node.(*parser.NodeLike), data)
	case reflect.TypeOf(new(parser.NodeNotLike)):
		return walkNotLike(node.(*parser.NodeNotLike), data)
	case reflect.TypeOf(new(parser.NodeNotEqual)):
		return walkNotEqual(node.(*parser.NodeNotEqual), data)
	case reflect.TypeOf(new(parser.NodeAnd)):
		return walkAnd(node.(*parser.NodeAnd), data)
	case reflect.TypeOf(new(parser.NodeOr)):
		return walkOr(node.(*parser.NodeOr), data)
	case reflect.TypeOf(new(parser.NodeGreater)):
		return walkGreater(node.(*parser.NodeGreater), data)
	case reflect.TypeOf(new(parser.NodeSmaller)):
		return walkSmaller(node.(*parser.NodeSmaller), data)
	default:
		return nil, fmt.Errorf("unknown node: %#v", node)
	}
}

func walkEqual(node *parser.NodeEqual, data *Problems) (*Problems, error) {
	return filter(data, node.LeftValue().(*parser.NodeId).Value(), node.RightValue().(*parser.NodeLiteral).Value(), true)
}

func walkLike(node *parser.NodeLike, data *Problems) (*Problems, error) {
	return filterLike(data, node.LeftValue().(*parser.NodeId).Value(), node.RightValue().(*parser.NodeLiteral).Value(), true)
}

func walkNotLike(node *parser.NodeNotLike, data *Problems) (*Problems, error) {
	return filterLike(data, node.LeftValue().(*parser.NodeId).Value(), node.RightValue().(*parser.NodeLiteral).Value(), false)
}

func walkNotEqual(node *parser.NodeNotEqual, data *Problems) (*Problems, error) {
	return filter(data, node.LeftValue().(*parser.NodeId).Value(), node.RightValue().(*parser.NodeLiteral).Value(), false)
}

func walkGreater(node *parser.NodeGreater, data *Problems) (*Problems, error) {
	return compare(data, node.LeftValue().(*parser.NodeId).Value(), node.RightValue().(*parser.NodeLiteral).Value(), false)
}

func walkSmaller(node *parser.NodeSmaller, data *Problems) (*Problems, error) {
	return compare(data, node.LeftValue().(*parser.NodeId).Value(), node.RightValue().(*parser.NodeLiteral).Value(), true)
}

func walkAnd(node *parser.NodeAnd, data *Problems) (*Problems, error) {
	// пересечение множеств в node.Left и node.Rigth
	left, err := walkNode(node.LeftValue(), data)
	if err != nil {
		return nil, err
	}
	rigth, err := walkNode(node.RightValue(), left) // в случае логических проблем тут должна находиться data
	if err != nil {
		return nil, err
	}
	return left.And(rigth), nil
}

func walkOr(node *parser.NodeOr, data *Problems) (*Problems, error) {
	// объединение множеств
	left, err := walkNode(node.LeftValue(), data)
	if err != nil {
		return nil, err
	}
	rigth, err := walkNode(node.RightValue(), data)
	if err != nil {
		return nil, err
	}
	return left.Or(rigth), nil
}

func compare(data *Problems, left, right string, bigger bool) (*Problems, error) {
	switch left {
	case "time":
		t, err := strconv.ParseInt(right, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse time: '%s'", err.Error())
		}
		result := NewProblems()
		for _, p := range data.Items {
			if bigger == (t > p.Time) {
				result.Add(p)
			}
		}
		return result, nil
	default:
		return nil, fmt.Errorf("unknown column for equal: '%s'", left)
	}
}

func filter(data *Problems, left, right string, eq bool) (*Problems, error) {
	switch left {
	case "state", "domain", "fqdn", "hostname", "host", "service":
		return filterText(data, left, right, eq), nil
	case "tag":
		return filterTag(data, right, eq), nil
	default:
		return nil, fmt.Errorf("unknown column for equal: '%s'", left)
	}
}

func filterText(data *Problems, left, right string, eq bool) *Problems {
	result := NewProblems()
	for _, p := range data.Items {
		if eq == (right == p.GetTextValue(left)) {
			result.Add(p)
		}
	}
	return result
}

func filterTag(data *Problems, tag string, eq bool) *Problems {
	result := NewProblems()
	for _, p := range data.Items {
		if eq == stringInTagSlice(tag, p.Tags) {
			result.Add(p)
		}
	}
	return result
}

func stringInTagSlice(tag string, data []string) bool {
	for _, d := range data {
		if d == tag {
			return true
		}
	}
	return false
}

// к нам приходит массив из проблем (data), в left хранится поле по которому надо фильтровать
// в rigth находится "%бла%бла%", по которому надо отфильтровать проблемы
func filterLike(data *Problems, left, right string, eq bool) (*Problems, error) {
	switch left {
	case "state", "domain", "fqdn", "hostname", "host", "service":
		if reg, err := makeLikeStringRegexp(right); err != nil {
			return nil, err
		} else {
			return filterLikeText(data, left, reg, eq), nil
		}
	default:
		return nil, fmt.Errorf("unknown column for like: '%s'", left)
	}
}

func filterLikeText(data *Problems, left string, right *regexp.Regexp, eq bool) *Problems {
	result := NewProblems()
	for _, problem := range data.Items {
		if eq == (right.MatchString(problem.GetTextValue(left))) {
			result.Add(problem)
		}
	}
	return result
}

func makeLikeStringRegexp(rigth string) (*regexp.Regexp, error) {
	find := strings.Replace(rigth, `.`, `\.`, -1)
	find = strings.Replace(find, `:`, `\:`, -1)
	find = strings.Replace(find, `%`, `(.+)?`, -1)
	return regexp.Compile(find)
}
