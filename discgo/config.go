// Created 2019-12-05 by NGnius

package main

import (
    "os"
    "io/ioutil"
    "encoding/json"
)

const (
  GlobalConfigPath = "config.json"
)

var (
  GlobalConfiguration = Configuration{
    path:GlobalConfigPath,
    LogPath:"discgo.log",
    Token:"token",
    Bot:true,
    RetryDelay:5*1000*1000,
  } // default configuration
)

func LoadGlobalConfigFile(path string) (error){
  // determine if the file exists
  _, statErr := os.Stat(path)
  if statErr != nil && os.IsNotExist(statErr) {
    // file does not exist; create default config file
    GlobalConfiguration.path = path
    GlobalConfiguration.Save()
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

type Configuration struct {
    LogPath string `json:"log"`
    Token string `json:"token"`
    Bot bool `json:"bot"`
    RetryDelay int `json:"retry-delay"`
    path string
}

func (conf *Configuration) Save() error {
    file, openErr := os.Create(conf.path)
    if openErr != nil {
      return openErr
    }
    defer file.Close()
    out, marshalErr := json.MarshalIndent(conf, "", " ")
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


