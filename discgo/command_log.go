// Created 2019-12-08 by NGnius

package main

import (
  "log"
  "os"
  //"strings"
  "github.com/bwmarrin/discordgo" // discord lib
)

func logCommandHandler(s *discordgo.Session, m * discordgo.MessageCreate) {
  if m.Author.ID == s.State.User.ID { // author == discgo
    return
  }
  if m.Content == "$log" {
    logpath, ok := GlobalConfiguration.GetValue("main", "log")
    if !ok {
      _, sendErr := s.ChannelMessageSend(m.ChannelID, "No log file configured")
      if sendErr != nil {
        log.Println("Failed to send log command failure message")
        log.Println(sendErr)
      }
      return
    }
    file, openErr := os.Open(logpath)
    if openErr != nil {
      _, sendErr := s.ChannelMessageSend(m.ChannelID, "Failed to open log file")
      if sendErr != nil {
        log.Println("Failed to send log command failure message")
        log.Println(sendErr)
      }
      return
    }
    uploads := []*discordgo.File{
      &discordgo.File{
        Name: "discgo.log",
        ContentType: "text/plain",
        Reader: file,
      },
    }
    complexMsg := discordgo.MessageSend {
      Content: "Plaintext log file",
      Files: uploads,
    }
    _, sendErr := s.ChannelMessageSendComplex(m.ChannelID, &complexMsg)
    if sendErr != nil {
      log.Println("Failed to send log command response message")
      log.Println(sendErr)
    }
  }
}
