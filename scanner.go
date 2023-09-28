package main

import (
  "fmt"
  "os/exec"
)

func scan_tcp(ip string) {
  cmd := exec.Command("nmap", "-sV", ip, "--stylesheet", "https://raw.githubusercontent.com/honze-net/nmap-bootstrap-xsl/master/nmap-bootstrap.xsl", "-oA", "scan", "-Pn")

  output, err := cmd.CombinedOutput()
  if err != nil {
      fmt.Println("Error executing command:", err)
      return
  }

  // Print the output of the command
  fmt.Println("Command Output:")
  fmt.Println(string(output))
}
