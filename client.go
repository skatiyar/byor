package byor

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	terminal "golang.org/x/term"
)

var clog = newLogger("[CLIENT]")

func Client(port string) error {
	addr, addrErr := net.ResolveTCPAddr("tcp", port)
	if addrErr != nil {
		return addrErr
	}

	tcpConn, connErr := net.DialTCP("tcp", nil, addr)
	if connErr != nil {
		return connErr
	}
	defer tcpConn.Close()

	oldState, termErr := terminal.MakeRaw(int(os.Stdin.Fd()))
	if termErr != nil {
		return termErr
	}
	defer terminal.Restore(int(os.Stdin.Fd()), oldState)

	screen := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}

	term := terminal.NewTerminal(screen, "")
	term.SetPrompt(string(term.Escape.Green) + "> " + string(term.Escape.Reset))

	welcomeMessage := formatOutput(term, "Connected to Server at address - "+port)
	fmt.Fprintln(term, welcomeMessage)
	cmdHelp := formatOutput(term, "To exit prompt use command (ctrl+c)")
	fmt.Fprintln(term, cmdHelp)

	for {
		line, rErr := term.ReadLine()
		if rErr != nil {
			if rErr == io.EOF {
				return nil
			}
			return rErr
		}

		line = strings.Trim(line, " ")
		if line == "" {
			continue
		}

		if wErr := Encoder(tcpConn, line); wErr != nil {
			fmt.Fprintln(term, formatErrorOutput(term, fmt.Sprintf("Writing to server: %v", wErr)))
			continue
		}

		data, dataErr := Decoder(tcpConn)
		if dataErr != nil {
			fmt.Fprintln(term, formatErrorOutput(term, fmt.Sprintf("Reading from server: %v", dataErr)))
			continue
		}
		fmt.Fprintln(term, formatServerReply(term, data))
	}
}

func formatOutput(term *terminal.Terminal, str string) string {
	return string(term.Escape.Cyan) + "# " + str + string(term.Escape.Reset)
}

func formatServerReply(term *terminal.Terminal, str string) string {
	return string(term.Escape.Yellow) + "[REPLY] " + str + string(term.Escape.Reset)
}

func formatErrorOutput(term *terminal.Terminal, err string) string {
	return string(term.Escape.Red) + "[ERROR] " + err + string(term.Escape.Reset)
}
