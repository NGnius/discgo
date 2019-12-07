// Created 2019-10-03 by NGnius

package main

import (
  "os"
  "log"
  "time"
  "github.com/bwmarrin/discordgo" // discord lib
)

var (
  disconnectChan = make(chan bool)
  session *discordgo.Session
  startupSuccessful bool
  connectSuccessful bool
  disconnecting bool
)

func init() {
  var err error
  // load configuration file
  err = LoadGlobalConfigFile(GlobalConfigPath)
  if err != nil {
    log.Println("Unable to read config file, what am I supposed to do?")
    log.Println(err)
    return
  }
  // configure log output file
  var f *os.File
  f, err = os.Create(GlobalConfiguration.LogPath)
  if err != nil {
    log.Println("Failed to create log file, attempting to continue")
    log.Println(err)
  } else {
    log.SetOutput(f) 
  }
  discordErr := initDiscord()
  if discordErr != nil {
    log.Println("Failed to initialize Discord session during startup, aborting")
    log.Println(discordErr)
    return
  }
  GlobalConfiguration.Save() // save out fixed config, in case it was formatted badly or missing values
  startupSuccessful = true
}

func main() {
  if startupSuccessful {
    log.Println("Startup successful")
    // register disconnect handler
    runLoop: for {
      shouldDisconnect := <- disconnectChan
      sessionDisconnectErr := session.Close()
      if sessionDisconnectErr != nil {
        log.Println("Unable to close Discord session")
        log.Println(sessionDisconnectErr)
      }
      disconnecting = false
      if shouldDisconnect {
        log.Println("Connection to Discord failed during runtime and instructed to not continue")
        break runLoop
      } else {
        // completely re-init discord
        time.Sleep(time.Duration(GlobalConfiguration.RetryDelay))
        discordErr := initDiscord()
        if discordErr != nil {
          log.Println("Failed to initialize Discord session during disconnect recovery, aborting")
          log.Println(discordErr)
          return
        }
        log.Println("Recreated Discord session after disconnect")
      }
    }
  } else {
    log.Println("Startup failed, no recovery possible")
    return
  }
  log.Println("DiscGo shutdown")
}

func initDiscord() error {
  // create discord session
  var err error
  session, err = discordgo.New(GlobalConfiguration.Token)
  if err != nil {
    return err
  }
  session.AddHandler(connectHandler)
  session.AddHandler(disconnectHandler)
  sessionErr := session.Open()
  if sessionErr != nil {
    return sessionErr
  }
  return nil
}

// discord event handlers for connection handling

func connectHandler(s *discordgo.Session, m * discordgo.Connect) {
  log.Println("Connection to Discord succeeded :)")
  connectSuccessful = true
}

func disconnectHandler (s *discordgo.Session, m *discordgo.Disconnect) {
  log.Println("Disconnected from Discord :(")
  if !disconnecting {
    disconnecting = true
    disconnectChan <- !connectSuccessful
    connectSuccessful = false
  } else {
    log.Println("Already in disconnecting state!")
  }
}
