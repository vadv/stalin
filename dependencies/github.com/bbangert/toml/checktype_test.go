package toml

import (
	gs "github.com/rafrombrc/gospec/src/gospec"
	"time"
)

func CheckTypeSpec(c gs.Context) {
	var err error

	var tomlBlob = `
ranking = ["Springsteen", "J Geils"]

[bands.Springsteen]
type = "ignore_this"
started = 1973
albums = ["Greetings", "WIESS", "Born to Run", "Darkness"]
not_albums = ["Greetings", "WIESS", "Born to Run", "Darkness"]

[bands.J Geils]
started = 1970
albums = ["The J. Geils Band", "Full House", "Blow Your Face Out"]
`

	type classics struct {
		Ranking []string
		Bands   map[string]Primitive
	}

	c.Specify("check mapping", func() {
		// Do the initial decode. Reflection is delayed on Primitive values.
		var music classics
		var md MetaData
		md, err = Decode(tomlBlob, &music)
		c.Assume(err, gs.IsNil)

		empty_ignore := map[string]interface{}{}
		err = CheckType(md.mapping, music, empty_ignore)
		c.Assume(err, gs.IsNil)
	})
}

func DecodeStrictSpec(c gs.Context) {

	// This blob when used with an empty_ignore blob
	var testSimple = `
age = 250
andrew = "gallant"
kait = "brady"
now = 1987-07-05T05:45:00Z 
yesOrNo = true
pi = 3.14
colors = [
	["red", "green", "blue"],
	["cyan", "magenta", "yellow", "black"],
]

[Annoying.Cats]
plato = "smelly"
cauchy = "stupido"
`

	var tomlBlob = `
# Some comments.
[alpha]
ip = "10.0.0.1"

	[alpha.config]
	Ports = [8001, 8002]
	Location = "Toronto"
	Created = 1987-07-05T05:45:00Z
`

	type serverConfig struct {
		Ports    []int
		Location string
		Created  time.Time
	}

	type server struct {
		IPAddress string       `toml:"ip"`
		Config    serverConfig `toml:"config"`
	}

	type kitties struct {
		Plato  string
		Cauchy string
	}

	type simple struct {
		Age      int
		Colors   [][]string
		Pi       float64
		YesOrNo  bool
		Now      time.Time
		Andrew   string
		Kait     string
		Annoying map[string]kitties
	}

	type servers map[string]server

	var config servers
	var val simple
	var err error

	empty_ignore := map[string]interface{}{}
	_, err = DecodeStrict(tomlBlob, &config, empty_ignore)
	c.Assume(err, gs.IsNil)

	_, err = DecodeStrict(testSimple, &val, empty_ignore)
	c.Assume(err, gs.IsNil)

	_, err = DecodeStrict(testBadArg, &val, empty_ignore)
	c.Assume(err.Error(), gs.Equals, "Configuration contains key [not_andrew] which doesn't exist in struct")

}

func DecodeStrictInterfaceSpec(c gs.Context) {
	// Check that we can safely decode into an empty interface
	// properly

	var tomlBlob = `
[MyMultiDecoder]
type = "MultiDecoder"
order = ["MyJsonDecoder", "MyProtobufDecoder"]

[MyMultiDecoder.delegates.MyJsonDecoder]
type = "JsonDecoder"
encoding_name = "JSON"

[MyMultiDecoder.delegates.MyProtobufDecoder]
type = "ProtobufDecoder"
encoding_name = "PROTOCOL_BUFFER"
`

	var err error
	var obj interface{}
	empty_ignore := map[string]interface{}{}
	_, err = DecodeStrict(tomlBlob, &obj, empty_ignore)
	c.Assume(err, gs.IsNil)

	actualObj := obj.(map[string]interface{})
	multidecoder := actualObj["MyMultiDecoder"].(map[string]interface{})
	c.Expect(multidecoder["type"], gs.Equals, "MultiDecoder")
	order := multidecoder["order"].([]interface{})

	d1 := order[0].(string)
	d2 := order[1].(string)
	c.Expect(d1, gs.Equals, "MyJsonDecoder")
	c.Expect(d2, gs.Equals, "MyProtobufDecoder")
	delegates := multidecoder["delegates"].(map[string]interface{})

	myjson := delegates["MyJsonDecoder"].(map[string]interface{})
	myproto := delegates["MyProtobufDecoder"].(map[string]interface{})

	c.Expect(myjson["type"], gs.Equals, "JsonDecoder")
	c.Expect(myjson["encoding_name"], gs.Equals, "JSON")

	c.Expect(myproto["type"], gs.Equals, "ProtobufDecoder")
	c.Expect(myproto["encoding_name"], gs.Equals, "PROTOCOL_BUFFER")

}
