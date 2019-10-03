// Created 2019-10-03 by NGnius

package main

import (
  "os"
  "log"
  "encoding/json"
  "io/ioutil"
  "github.com/bwmarrin/discordgo"
)

const (
  configPath = "config.json"
)

var (
  disconnectChan = make(chan bool)
  configuration map[string]string
  session *discordgo.Session
  startupSuccessful bool
)

func init() {
  var err error
  // load configuration file
  err = loadConfig(configPath)
  if err != nil {
    log.Println(err)
    return
  }
  // configure log output file
  var f *os.File
  f, err = os.Create(configuration["log"])
  if err != nil {
    log.Println(err)
    return
  }
  log.SetOutput(f)
  // create discord session
  session, err = discordgo.New(configuration["token"])
  if err != nil {
    log.Println(err)
    return
  }
  startupSuccessful = true
}

func main() {
  if startupSuccessful {
    log.Println("Startup successful")
    log.Println("Opening session")
    sessionErr := session.Open()
    if sessionErr != nil {
      log.Println(sessionErr)
      return
    }
    // register disconnect handler
    session.AddHandler(func(s *discordgo.Session, m *discordgo.Disconnect) {
      disconnectChan <- true
    })
    runLoop: for {
      isDisconnected := <- disconnectChan
      if !isDisconnected {
        break runLoop
      }
      sessionErr := session.Open()
      if sessionErr != nil {
        log.Println(sessionErr)
        return
      }
    }
  } else {
    log.Println("Startup failed, please read the log for details")
    return
  }
}

func loadConfig(path string) (error){
  file, openErr := os.Open(path)
  defer file.Close()
  if openErr != nil {
    return openErr
  }
  data, readErr := ioutil.ReadAll(file)
  if readErr != nil {
    return readErr
  }
  unmarshalErr := json.Unmarshal(data, &configuration)
  if unmarshalErr != nil {
    return unmarshalErr
  }
  return nil
}
