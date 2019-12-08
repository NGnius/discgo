// Created 2019-12-05 by NGnius

package main

import (
    "os"
    "io/ioutil"
    "encoding/json"
)

const (
  DefaultGlobalConfigPath = "config.json"
)

var (
  DefaultConfiguration = MasterConfiguration{
    path: DefaultGlobalConfigPath,
    Version: "0000",
    Configurations: map[string]SlaveConfiguration{
      "main": SlaveConfiguration{
        Name: "main",
        Description: "Discgo's main configuration information (required)",
        Mappings: map[string]string {
          "log": "discgo.log",
          "debug-channelID": "channelID",
          "token": "token",
          "retry-delay": "5s",
        },
      },
    },
  }
  GlobalConfiguration = DefaultConfiguration
)

func LoadGlobalConfigFile(path string) (error) {
  // determine if the file exists
  _, statErr := os.Stat(path)
  if statErr != nil && os.IsNotExist(statErr) {
    // file does not exist; use default config
    GlobalConfiguration.path = path
    return nil
  }
  file, openErr := os.Open(path)
  defer file.Close()
  if openErr != nil {
    return openErr
  }
  data, readErr := ioutil.ReadAll(file)
  if readErr != nil {
    return readErr
  }
  unmarshalErr := json.Unmarshal(data, &GlobalConfiguration)
  if unmarshalErr != nil {
    return unmarshalErr
  }
  GlobalConfiguration.path = path
  return nil
}

func LoadDefaultGlobalConfigFile() (error) {
  return LoadGlobalConfigFile(DefaultGlobalConfigPath)
}

type MasterConfiguration struct {
    Version string `json:"version"`
    Configurations map[string]SlaveConfiguration `json:"configurations"`
    path string
}

func (conf *MasterConfiguration) Save() error {
    file, openErr := os.Create(conf.path)
    if openErr != nil {
      return openErr
    }
    defer file.Close()
    out, marshalErr := json.MarshalIndent(conf, "", "  ")
    if marshalErr != nil {
      return marshalErr
    }
    _, writeErr := file.Write(out)
    if writeErr != nil {
      return writeErr
    }
    syncErr := file.Sync()
    if syncErr != nil {
      return syncErr
    }
    return nil
}

func (conf *MasterConfiguration) GetValue(confKey string, mapKey string) (value string, ok bool) {
    var slave SlaveConfiguration
    slave, ok = conf.Configurations[confKey]
    if !ok {
        return
    }
    return slave.GetValue(mapKey)
}

func (conf *MasterConfiguration) TryGetValue(confKey string, mapKey string) (value string) {
    slave, ok := conf.Configurations[confKey]
    if !ok {
        return
    }
    return slave.TryGetValue(mapKey)
}

type SlaveConfiguration struct {
    Name string `json:"name"` // the name of the SlaveConfiguration, usually the key from MasterConfiguration.Mappings
    Description string `json:"description"` // a brief description of the use of the SlaveConfiguration
    Mappings map[string]string `json:"mappings"` // key:value pairs
}

func (slaveConf *SlaveConfiguration) GetValue(key string) (value string, ok bool) {
    value, ok = slaveConf.Mappings[key]
    return
}

func (slaveConf *SlaveConfiguration) TryGetValue(key string) (value string) {
    value = slaveConf.Mappings[key]
    return
}
