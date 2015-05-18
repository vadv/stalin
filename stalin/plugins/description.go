package plugins

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"text/tabwriter"
)

// build plugin description by tag "description"
func getDescription(p pluginCreator) string {

	plugin := p.Creator("description")
	typ := reflect.TypeOf(plugin)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	val := reflect.ValueOf(plugin)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	writer := &tabwriter.Writer{}
	doc := &bytes.Buffer{}
	writer.Init(doc, 0, 8, 0, '\t', tabwriter.AlignRight)

	fmt.Fprintf(writer, "Type: '%s'\n", p.Type)
	fmt.Fprintf(writer, "Description: %v\n", p.Description)
	if isOutputPlugin(plugin) {
		fmt.Fprintf(writer, "Access to external messages: yes\n")
	} else {
		fmt.Fprintf(writer, "Access to external messages: no\n")
	}
	fmt.Fprintf(writer, "Config:\n")

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tags := field.Tag
		if !field.Anonymous {
			if tag := tags.Get("description"); tag != "" {
				//method := reflect.
				name := ""
				if tags.Get("json") != "" {
					jsons := strings.Split(tags.Get("json"), ",")
					name = jsons[0]
				} else {
					name = field.Name
				}
				value := val.FieldByName(field.Name).Interface()
				fmt.Fprintf(writer, "\t* '%s' (%v, %#v)\t'%s'\n", name, field.Type, value, tag)
			}
		}
	}

	writer.Flush()
	return doc.String()
}

func PluginDescription(pluginType string) string {
	if pluginType == "all" {
		return allPluginDescription()
	}
	if creator, ok := availablePlugins[pluginType]; !ok {
		return fmt.Sprintf("Plugin type '%s' not found. List of aviable plugins:\n%s", pluginType, allPluginDescription())
	} else {
		return getDescription(creator)
	}
}

func allPluginDescription() (description string) {
	for _, typ := range availablePlugins.getSortedTypes() {
		description = fmt.Sprintf("%s\n\n%s", description, getDescription(availablePlugins[typ]))
	}
	return
}
