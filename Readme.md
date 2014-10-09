# Stalin

## Use cases
```
Input (http)  -         - Output (carbon)
Input (tcp)   -  Router - Output (over stalin)
Input (etc..) -         - Output (pg, etc...)
````

## Example config
```toml
[input_tcp]
type = "TcpInput"
address = "0.0.0.0:5555"

[output_tcp_first]
type = "TcpOutput"
address = "direction1.host:5555"

[output_tcp_second]
type = "TcpOutput"
address = "direction2.host:5555"

# some format to input json
[input_http]
type = "HttpInput"
address = "0.0.0.0:80"
queue = "MessagesJson1"

# output to elastic
[input_http]
type = "ElasticOut"
address = "elast.host:80"
queue = "MessagesJson1"
```

## Example plugin
```go
package plugin_file

import (
  "encoding/json"
  "fmt"
  "stalin/message"
  . "stalin/plugins"
)

func init() {
  RegisterPlugin("PluginFileOutput", new(PluginFileOutput))
}

type PluginFileOutput struct {
  gConfig *GlobalConfig
  config  *PluginFileOutputConfig
}

type PluginFileOutputConfig struct {
  Path         string `json:"path"` // values in config file
}

func (r *PluginFileOutput) Init(g *GlobalConfig, name string) (Plugin, error) {
  plugin := &PluginFileOutputConfig{gConfig: g}
  config := &PluginFileOutputConfig{
    Host:         "/tmp/file.log", // default values
  }
  if err := json.Unmarshal(g.PluginConfigs[name], config); err != nil {
    return nil, err
  }
  plugin.config = config
  return plugin, nil
}

func (r *PluginFileOutput) Inject(msg *message.Message) error {
// event for new message
  return nil
}

func (r *PluginFileOutput) Run() error {
  return nil
}
```
