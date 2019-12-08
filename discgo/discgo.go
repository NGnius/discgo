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
  Session *discordgo.Session
  startupSuccessful bool
  connectSuccessful bool
  disconnecting bool
)

func init() {
  var err error
  // load configuration file
  err = LoadDefaultGlobalConfigFile()
  if err != nil {
    log.Println("Unable to read config file, what am I supposed to do?")
    log.Println(err)
    return
  }
  GlobalConfiguration.Save() // save out fixed config, in case it was old or missing
  // configure log output file
  var f *os.File
  f, err = os.Create(GlobalConfiguration.Configurations["main"].Mappings["log"])
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
  startupSuccessful = true
}

func main() {
  if startupSuccessful {
    log.Println("Startup successful")
    // register disconnect handler
    runLoop: for {
      shouldDisconnect := <- disconnectChan
      Session.Logout()
      sessionDisconnectErr := Session.Close()
      if sessionDisconnectErr != nil {
        log.Println("Unable to close Discord session")
        log.Println(sessionDisconnectErr)
      }
      if shouldDisconnect {
        log.Println("Connection to Discord failed during runtime and instructed to not continue")
        break runLoop
      } else {
        // completely re-init discord
        dur, parseErr := time.ParseDuration(GlobalConfiguration.Configurations["main"].Mappings["retry-delay"])
        if parseErr != nil {
          log.Println("Unable to parse main::retry-delay config value, aborting")
          return
        }
        time.Sleep(dur)
        disconnecting = false
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
  connectSuccessful = false
  var err error
  Session, err = discordgo.New("Bot "+GlobalConfiguration.Configurations["main"].Mappings["token"])
  if err != nil {
    return err
  }
  Session.AddHandler(connectHandler)
  Session.AddHandler(disconnectHandler)
  sessionErr := Session.Open()
  if sessionErr != nil {
    return sessionErr
  }
  registerMessageHandlers(Session)
  return nil
}

func registerMessageHandlers(s *discordgo.Session) {
  // register command message handlers here
  s.AddHandler(logCommandHandler)
}

// discord event handlers for connection handling

func connectHandler(s *discordgo.Session, m * discordgo.Connect) {
  log.Println("Connection to Discord succeeded :)")
  sendErr := sendDebugMessage("Connection to Discord succeeded :)")
  if sendErr != nil {
    log.Println("Failed to send Discord message on connect")
    log.Println(sendErr)
  }
  connectSuccessful = true
}

func disconnectHandler (s *discordgo.Session, m *discordgo.Disconnect) {
  log.Println("Disconnected from Discord :(")
  if !disconnecting {
    disconnecting = true
    connectSuccessful = false
    disconnectChan <- false // always try to reconnect
  } else {
    log.Println("Already in disconnecting state!")
  }
}

// discord tools & macros

func sendDebugMessage(content string) (err error) {
  channelID, ok := GlobalConfiguration.GetValue("main", "debug-channelID")
  if (!ok || len(channelID) != 18) {
    return // fail silently if not valid channel ID
  }
  _, err = Session.ChannelMessageSend(channelID, content)
  return
}
