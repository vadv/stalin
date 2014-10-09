package toml

import (
	gs "github.com/rafrombrc/gospec/src/gospec"
	"log"
	"reflect"
	"strings"
	"time"
)

func init() {
	log.SetFlags(0)
}

var testBadArg = `
age = 250

not_andrew = "gallant"
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

func DecodeSpec(c gs.Context) {
	var val simple

	md, err := Decode(testSimple, &val)
	c.Assume(err, gs.IsNil)
	c.Assume(md.IsDefined("Annoying", "Cats", "plato"), gs.IsTrue)
	c.Assume(md.IsDefined("Cats", "Stinky"), gs.IsFalse)
	var colors = [][]string{[]string{"red", "green", "blue"},
		[]string{"cyan", "magenta", "yellow", "black"}}
	for ridx, row := range colors {
		for cidx, _ := range row {
			c.Assume(val.Colors[ridx][cidx], gs.Equals, colors[ridx][cidx])
		}
	}
	c.Assume(val, gs.Not(gs.IsNil))
}

// Case insensitive matching tests.
// A bit more comprehensive than needed given the current implementation,
// but implementations change.
// Probably still missing demonstrations of some ugly corner cases regarding
// case insensitive matching and multiple fields.
var caseToml = `
tOpString = "string"
tOpInt = 1
tOpFloat = 1.1
tOpBool = true
tOpdate = 2006-01-02T15:04:05Z
tOparray = [ "array" ]
Match = "i should be in Match only"
MatcH = "i should be in MatcH only"
Field = "neat"
FielD = "messy"
once = "just once"
[nEst.eD]
nEstedString = "another string"
`

type Insensitive struct {
	TopString string
	TopInt    int
	TopFloat  float64
	TopBool   bool
	TopDate   time.Time
	TopArray  []string
	Match     string
	MatcH     string
	Field     string
	Once      string
	OncE      string
	Nest      InsensitiveNest
}

type InsensitiveNest struct {
	Ed InsensitiveEd
}

type InsensitiveEd struct {
	NestedString string
}

func InsensitiveDecodeSpec(c gs.Context) {
	tme, err := time.Parse(time.RFC3339, time.RFC3339[:len(time.RFC3339)-5])
	if err != nil {
		panic(err)
	}
	expected := Insensitive{
		TopString: "string",
		TopInt:    1,
		TopFloat:  1.1,
		TopBool:   true,
		TopDate:   tme,
		TopArray:  []string{"array"},
		MatcH:     "i should be in MatcH only",
		Match:     "i should be in Match only",
		Field:     "neat", // encoding/json would store "messy" here
		Once:      "just once",
		OncE:      "just once", // wait, what?
		Nest: InsensitiveNest{
			Ed: InsensitiveEd{NestedString: "another string"},
		},
	}
	var got Insensitive
	_, err = Decode(caseToml, &got)
	c.Assume(err, gs.IsNil)
	c.Assume(reflect.DeepEqual(expected, got), gs.IsTrue)
}

func ExamplePrimitiveDecodeSpec(c gs.Context) {
	var md MetaData
	var err error

	var tomlBlob = `
ranking = ["Springsteen", "J Geils"]

[bands.Springsteen]
started = 1973
albums = ["Greetings", "WIESS", "Born to Run", "Darkness"]

[bands.J Geils]
started = 1970
albums = ["The J. Geils Band", "Full House", "Blow Your Face Out"]
`

	type band struct {
		Started int
		Albums  []string
	}

	type classics struct {
		Ranking []string
		Bands   map[string]Primitive
	}

	// Do the initial decode. Reflection is delayed on Primitive values.
	var music classics
	if md, err = Decode(tomlBlob, &music); err != nil {
		log.Fatal(err)
	}

	c.Specify("decode with pointers to struct", func() {
		// MetaData still includes information on Primitive values.
		c.Assume(md.IsDefined("bands", "Springsteen"), gs.IsTrue)

		// Decode primitive data into Go values.
		expected_artists := map[string]int{"Springsteen": 1973, "J Geils": 1970}

		for _, artist := range music.Ranking {
			// A band is a primitive value, so we need to decode it to get a
			// real `band` value.
			primValue := music.Bands[artist]

			var aBand band
			err = PrimitiveDecode(primValue, &aBand)
			c.Assume(err, gs.IsNil)
			c.Assume(expected_artists[artist], gs.Equals, aBand.Started)
		}
	})

	c.Specify("decode with struct fails", func() {
		for _, artist := range music.Ranking {
			// A band is a primitive value, so we need to decode it to get a
			// real `band` value.
			primValue := music.Bands[artist]

			var aBand band
			err = PrimitiveDecode(primValue, aBand)
			c.Assume(err, gs.Not(gs.IsNil))
			c.Expect(strings.HasPrefix(err.Error(), "Can't use non-pointer"), gs.IsTrue)
		}
	})
}

func ExampleDecodeSpec(c gs.Context) {
	var tomlBlob = `
	# Some comments.
	[alpha]
	ip = "10.0.0.1"

		[alpha.config]
		Ports = [ 8001, 8002 ]
		Location = "Toronto"
		Created = 1987-07-05T05:45:00Z

	[beta]
	ip = "10.0.0.2"

		[beta.config]
		Ports = [ 9001, 9002 ]
		Location = "New Jersey"
		Created = 1887-01-05T05:55:00Z
	`

	type serverConfig struct {
		Ports    []int
		Location string
		Created  time.Time
	}

	type server struct {
		IP     string       `toml:"ip"`
		Config serverConfig `toml:"config"`
	}

	type servers map[string]server

	var config servers
	if _, err := Decode(tomlBlob, &config); err != nil {
		log.Fatal(err)
	}

	const date_form = "2006-01-02T15:04:05Z"
	alpha_time, _ := time.Parse(date_form, "1987-07-05T05:45:00Z")
	beta_time, _ := time.Parse(date_form, "1887-01-05T05:55:00Z")

	expected_servers := map[string]interface{}{
		"alpha": map[string]interface{}{"ip": "10.0.0.1",
			"config": &serverConfig{Ports: []int{8001, 8002},
				Location: "Toronto",
				Created:  alpha_time,
			},
		},
		"beta": map[string]interface{}{"ip": "10.0.0.2",
			"config": &serverConfig{Ports: []int{9001, 9002},
				Location: "New Jersey",
				Created:  beta_time,
			},
		}}

	for _, name := range []string{"alpha", "beta"} {
		actual_node := config[name]

		expected_node := expected_servers[name].(map[string]interface{})

		actual_config := actual_node.Config
		expected_config := expected_node["config"].(*serverConfig)

		c.Expect(actual_node.IP, gs.Equals, expected_node["ip"])

		c.Expect(actual_config.Location, gs.Equals, expected_config.Location)
		c.Expect(actual_config.Created, gs.Equals, expected_config.Created)

		c.Expect(len(actual_config.Ports), gs.Equals, len(expected_config.Ports))
		for idx, _ := range actual_config.Ports {
			c.Expect(actual_config.Ports[idx], gs.Equals, expected_config.Ports[idx])
		}
	}

	// Server: alpha (ip: 10.0.0.1) in Toronto created on 1987-07-05
	// Ports: [8001 8002]
	// Server: beta (ip: 10.0.0.2) in New Jersey created on 1887-01-05
	// Ports: [9001 9002]
}

func PrimitiveDecodeStrictSpec(c gs.Context) {
	var md MetaData
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

	type band struct {
		Started int
		Albums  []string
	}

	type classics struct {
		Ranking []string
		Bands   map[string]Primitive
	}

	// Do the initial decode. Reflection is delayed on Primitive values.
	var music classics
	md, err = Decode(tomlBlob, &music)
	c.Assume(err, gs.IsNil)

	// MetaData still includes information on Primitive values.
	c.Assume(md.IsDefined("bands", "Springsteen"), gs.IsTrue)

	ignore_type := map[string]interface{}{"type": true}
	// Decode primitive data into Go values.
	for _, artist := range music.Ranking {
		// A band is a primitive value, so we need to decode it to get a
		// real `band` value.
		primValue := music.Bands[artist]

		var aBand band

		err = PrimitiveDecodeStrict(primValue, &aBand, ignore_type)
		if artist == "Springsteen" {
			c.Assume(err, gs.Not(gs.IsNil))
			c.Expect(err.Error(), gs.Equals, "Configuration contains key [not_albums] which doesn't exist in struct")
			c.Assume(1973, gs.Equals, aBand.Started)
		} else {
			c.Expect(err, gs.IsNil)
			c.Assume(1970, gs.Equals, aBand.Started)
		}

	}
}
