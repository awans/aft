package db

import (
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/olebedev/go-duktape.v3"
	"net/url"
	"regexp"
)

var (
	ErrValue = fmt.Errorf("%w: invalid value for type", ErrData)
)

// Id is the GUID of the datatype
// Name is the plain english name of the type
// FromJson takes as input the text to be parsed, returning it if it's valid
// Type returns the type of the raw data stored for the datatype
type Datatype struct {
	Id       uuid.UUID
	Name     string
	FromJson func(interface{}) (interface{}, error)
	Type     interface{}
}

var Bool = Datatype{
	Id:       uuid.MustParse("ca05e233-b8a2-4c83-a5c8-87b461c87184"),
	Name:     "bool",
	FromJson: boolFromJson,
	Type:     false,
}
var Int = Datatype{
	Id:       uuid.MustParse("17cfaaec-7a75-4035-8554-83d8d9194e97"),
	Name:     "int",
	FromJson: intFromJson,
	Type:     int64(0),
}
var Enum = Datatype{
	Id:       uuid.MustParse("f9e66ef9-2fa3-4588-81c1-b7be6a28352e"),
	Name:     "enum",
	FromJson: enumFromJson,
	Type:     int64(0),
}
var String = Datatype{
	Id:       uuid.MustParse("cbab8b98-7ec3-4237-b3e1-eb8bf1112c12"),
	Name:     "string",
	FromJson: stringFromJson,
	Type:     "",
}
var Text = Datatype{
	Id:       uuid.MustParse("4b601851-421d-4633-8a68-7fefea041361"),
	Name:     "text",
	FromJson: textFromJson,
	Type:     "",
}
var EmailAddress = Datatype{
	Id:       uuid.MustParse("6c5e513b-9965-4463-931f-dd29751f5ae1"),
	Name:     "email address",
	FromJson: emailAddressFromJson,
	Type:     "",
}
var UUID = Datatype{
	Id:       uuid.MustParse("9853fd78-55e6-4dd9-acb9-e04d835eaa42"),
	Name:     "uuid",
	FromJson: uuidFromJson,
	Type:     uuid.UUID{},
}
var Float = Datatype{
	Id:       uuid.MustParse("72e095f3-d285-47e6-8554-75691c0145e3"),
	Name:     "float",
	FromJson: floatFromJson,
	Type:     0.0,
}
var URL = Datatype{
	Id:       uuid.MustParse("84c8c2c5-ff1a-4599-9605-b56134417dd7"),
	Name:     "url",
	FromJson: urlFromJson,
	Type:     "",
}
var NativeCode = Datatype{
	Id:       uuid.MustParse("f34e5dd5-9209-4ce0-81ef-8e2d1ee86ece"),
	Name:     "native code",
	FromJson: nativeCodeFromJson,
	Type:     "",
}
var Javascript = Datatype{
	Id:       uuid.MustParse("210fd74f-bb52-4ff6-b45a-d9fcd8d980ae"),
	Name:     "javascript",
	FromJson: javascriptFromJson,
	Type:     "",
}

var datatypeMap map[uuid.UUID]Datatype = map[uuid.UUID]Datatype{
	Bool.Id:         Bool,
	Int.Id:          Int,
	Enum.Id:         Enum,
	String.Id:       String,
	Text.Id:         Text,
	EmailAddress.Id: EmailAddress,
	UUID.Id:         UUID,
	Float.Id:        Float,
	URL.Id:          URL,
	NativeCode.Id:   NativeCode,
	Javascript.Id:   Javascript,
}

// bool datatype
func boolFromJson(value interface{}) (interface{}, error) {
	b, ok := value.(bool)
	if !ok {
		return nil, fmt.Errorf("%w: expected bool got %T", ErrValue, value)
	}
	return b, nil
}

// int datatype
func intFromJson(value interface{}) (interface{}, error) {
	f, ok := value.(float64)
	if ok {
		i := int64(f)
		if !ok {
			return nil, fmt.Errorf("%w: expected int/enum got %T", ErrValue, value)
		}
		return i, nil
	}
	intVal, ok := value.(int)
	if ok {
		i := int64(intVal)
		if !ok {
			return nil, fmt.Errorf("%w: expected int/enum got %T", ErrValue, value)
		}
		return i, nil
	}
	i64Val, ok := value.(int64)
	if ok {
		return i64Val, nil
	} else {
		return nil, fmt.Errorf("%w: expected int/enum got %T", ErrValue, value)
	}
}

// enum datatype
func enumFromJson(value interface{}) (interface{}, error) {
	f, ok := value.(float64)
	if ok {
		i := int64(f)
		if !ok {
			return nil, fmt.Errorf("%w: expected int/enum got %T", ErrValue, value)
		}
		return i, nil
	}
	intVal, ok := value.(int)
	if ok {
		i := int64(intVal)
		if !ok {
			return nil, fmt.Errorf("%w: expected int/enum got %T", ErrValue, value)
		}
		return i, nil
	}
	i64Val, ok := value.(int64)
	if ok {
		return i64Val, nil
	} else {
		return nil, fmt.Errorf("%w: expected int/enum got %T", ErrValue, value)
	}
}

// string datatype
func stringFromJson(value interface{}) (interface{}, error) {
	s, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("%w: expected string got %T", ErrValue, value)
	}
	return s, nil
}

func stringType() interface{} {
	return ""
}

// text datatype
func textFromJson(value interface{}) (interface{}, error) {
	s, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("%w: expected text got %T", ErrValue, value)
	}
	return s, nil
}

// Email Address datatype.
// This is an example of a more complex datatype
//https://www.alexedwards.net/blog/validation-snippets-for-go#email-validation
var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func matchEmail(s string) bool {
	return rxEmail.MatchString(s)
}

func emailAddressFromJson(value interface{}) (interface{}, error) {
	emailAddressString, ok := value.(string)
	if ok {
		if (len(emailAddressString) > 254 || !matchEmail(emailAddressString)) && len(emailAddressString) != 0 {
			return nil, fmt.Errorf("%w: expected email address got %v", ErrValue, emailAddressString)
		}
	} else {
		return nil, fmt.Errorf("%w: expected email address got %T", ErrValue, value)
	}
	return emailAddressString, nil
}

// UUID datatype. Uses UUID from google underneath
func uuidFromJson(value interface{}) (interface{}, error) {
	var u uuid.UUID
	uuidString, ok := value.(string)
	if ok {
		var err error
		u, err = uuid.Parse(uuidString)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrValue, err)
		}

	} else {
		u, ok = value.(uuid.UUID)
		if !ok {
			return nil, fmt.Errorf("%w: expected uuid got %T", ErrValue, value)
		}
	}
	return u, nil
}

// float datatype.
func floatFromJson(value interface{}) (interface{}, error) {
	f, ok := value.(float64)
	if !ok {
		return nil, fmt.Errorf("%w: expected float got %T", ErrValue, value)
	}
	return f, nil
}

// URL datatype
// This is an example of a more complex datatype.
func urlFromJson(value interface{}) (interface{}, error) {
	urlString, ok := value.(string)
	if ok {
		u, err := url.Parse(urlString)
		if err != nil {
			return nil, fmt.Errorf("%w: expected url got %T", ErrValue, value)
		} else if u.Scheme == "" || u.Host == "" {
			return nil, fmt.Errorf("%w: expected url got %T", ErrValue, value)
		} else if u.Scheme != "http" && u.Scheme != "https" {
			return nil, fmt.Errorf("%w: expected url got %T", ErrValue, value)
		}
	} else {
		return nil, fmt.Errorf("%w: expected url got %T", ErrValue, value)
	}
	return urlString, nil
}

// Native code datatype
// This shouldn't ever get called, but is a logical placeholder for the database
func nativeCodeFromJson(value interface{}) (interface{}, error) {
	return nil, fmt.Errorf("%w: nativeCode does not execute type", ErrValue)
}

// Javascript datatype
// uses https://github.com/olebedev/go-duktape bindings
func javascriptFromJson(value interface{}) (interface{}, error) {
	javascriptString, ok := value.(string)
	if ok {
		ctx := duktape.New()
		err := ctx.PevalString(javascriptString)
		result := ctx.GetNumber(-1)
		ctx.DestroyHeap()
		if &err != nil {
			return result, nil
		} else {
			return nil, fmt.Errorf("%w: expected Javascript got %s", ErrValue, err)
		}
	}
	return nil, fmt.Errorf("%w: expected Javascript got %s", ErrValue, value)
}