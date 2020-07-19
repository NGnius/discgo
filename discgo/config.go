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
  DefaultConfiguration = MainConfiguration{
    path: DefaultGlobalConfigPath,
    Version: "0000",
    Configurations: map[string]SubordinateConfiguration{
      "main": SubordinateConfiguration{
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

type MainConfiguration struct {
    Version string `json:"version"`
    Configurations map[string]SubordinateConfiguration `json:"configurations"`
    path string
}

func (conf *MainConfiguration) Save() error {
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

func (conf *MainConfiguration) GetValue(confKey string, mapKey string) (value string, ok bool) {
    var subordinate SubordinateConfiguration
    subordinate, ok = conf.Configurations[confKey]
    if !ok {
        return
    }
    return subordinate.GetValue(mapKey)
}

func (conf *MainConfiguration) TryGetValue(confKey string, mapKey string) (value string) {
    subordinate, ok := conf.Configurations[confKey]
    if !ok {
        return
    }
    return subordinate.TryGetValue(mapKey)
}

type SubordinateConfiguration struct {
    Name string `json:"name"` // the name of the SubordinateConfiguration, usually the key from MainConfiguration.Mappings
    Description string `json:"description"` // a brief description of the use of the SubordinateConfiguration
    Mappings map[string]string `json:"mappings"` // key:value pairs
}

func (subordinateConf *SubordinateConfiguration) GetValue(key string) (value string, ok bool) {
    value, ok = subordinateConf.Mappings[key]
    return
}

func (subordinateConf *SubordinateConfiguration) TryGetValue(key string) (value string) {
    value = subordinateConf.Mappings[key]
    return
}
