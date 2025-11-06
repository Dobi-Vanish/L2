package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

// Shell contains settings for shell.
type Shell struct {
	reader    *bufio.Reader
	writer    *bufio.Writer
	running   bool
	lastError error
}

func NewShell() *Shell {
	return &Shell{
		reader: bufio.NewReader(os.Stdin),
		writer: bufio.NewWriter(os.Stdout),
	}
}

func (s *Shell) Run() {
	s.running = true

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for sig := range signalChan {
			if sig == syscall.SIGINT {
				fmt.Fprintln(s.writer, "\nReceived SIGINT")
				s.writer.Flush()
			}
		}
	}()

	for s.running {
		s.printPrompt()
		line, err := s.readInput()
		if err != nil {
			if err == io.EOF {
				fmt.Fprintln(s.writer, "\nExiting...")
				break
			}
			fmt.Fprintf(s.writer, "Error reading input: %v\n", err)
			continue
		}

		if strings.TrimSpace(line) == "" {
			continue
		}

		s.executeLine(line)
	}
}

func (s *Shell) printPrompt() {
	wd, err := os.Getwd()
	if err != nil {
		wd = "?"
	}
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "localhost"
	}

	fmt.Fprintf(s.writer, "\033[32m%s@%s\033[0m:\033[34m%s\033[0m$ ", os.Getenv("USER"), hostname, filepath.Base(wd))
	s.writer.Flush()
}

func (s *Shell) readInput() (string, error) {
	return s.reader.ReadString('\n')
}

func (s *Shell) executeLine(line string) {
	line = strings.TrimSpace(line)

	commands := strings.Split(line, "|")

	if len(commands) == 1 {
		s.executeSingleCommand(commands[0])
	} else {
		s.executePipeline(commands)
	}
}

func (s *Shell) executeSingleCommand(cmdLine string) {
	cmdLine = strings.TrimSpace(cmdLine)

	if s.executeBuiltinCommand(cmdLine) {
		return
	}

	args := s.parseArguments(cmdLine)
	if len(args) == 0 {
		return
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)
	defer signal.Stop(signalChan)

	go func() {
		<-signalChan
		if cmd.Process != nil {
			cmd.Process.Signal(syscall.SIGINT)
		}
	}()

	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(s.writer, "Error executing command: %v\n", err)
		s.lastError = err
	} else {
		s.lastError = nil
	}
	s.writer.Flush()
}

func (s *Shell) executeBuiltinCommand(cmdLine string) bool {
	args := s.parseArguments(cmdLine)
	if len(args) == 0 {
		return false
	}

	switch args[0] {
	case "cd":
		s.cdCommand(args[1:])
		return true
	case "pwd":
		s.pwdCommand(args[1:])
		return true
	case "echo":
		s.echoCommand(args[1:])
		return true
	case "kill":
		s.killCommand(args[1:])
		return true
	case "ps":
		s.psCommand(args[1:])
		return true
	case "exit":
		s.running = false
		return true
	default:
		return false
	}
}

func (s *Shell) cdCommand(args []string) {
	var path string
	if len(args) == 0 {
		path = os.Getenv("HOME")
		if path == "" {
			fmt.Fprintln(s.writer, "cd: HOME not set")
			return
		}
	} else {
		path = args[0]
	}

	err := os.Chdir(path)
	if err != nil {
		fmt.Fprintf(s.writer, "cd: %v\n", err)
		s.lastError = err
	} else {
		s.lastError = nil
	}
}

func (s *Shell) pwdCommand(args []string) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(s.writer, "pwd: %v\n", err)
		s.lastError = err
	} else {
		fmt.Fprintln(s.writer, wd)
		s.lastError = nil
	}
	s.writer.Flush()
}

func (s *Shell) echoCommand(args []string) {
	for i, arg := range args {
		if i > 0 {
			fmt.Fprint(s.writer, " ")
		}
		fmt.Fprint(s.writer, arg)
	}
	fmt.Fprintln(s.writer)
	s.writer.Flush()
	s.lastError = nil
}

func (s *Shell) killCommand(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(s.writer, "kill: usage: kill <pid>")
		s.lastError = errors.New("missing pid")
		return
	}

	pid, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintf(s.writer, "kill: invalid pid: %v\n", err)
		s.lastError = err
		return
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Fprintf(s.writer, "kill: %v\n", err)
		s.lastError = err
		return
	}

	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		fmt.Fprintf(s.writer, "kill: %v\n", err)
		s.lastError = err
	} else {
		fmt.Fprintf(s.writer, "Sent SIGTERM to process %d\n", pid)
		s.lastError = nil
	}
	s.writer.Flush()
}

func (s *Shell) psCommand(args []string) {
	cmd := exec.Command("ps", "aux")
	cmd.Stdout = s.writer
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(s.writer, "ps: %v\n", err)
		s.lastError = err
	} else {
		s.lastError = nil
	}
	s.writer.Flush()
}

func (s *Shell) executePipeline(commands []string) {
	var cmds []*exec.Cmd
	var err error

	for _, cmdStr := range commands {
		args := s.parseArguments(strings.TrimSpace(cmdStr))
		if len(args) == 0 {
			continue
		}

		var cmd *exec.Cmd
		if s.executeBuiltinCommand(cmdStr) {
			fmt.Fprintln(s.writer, "Builtin commands not supported in pipelines")
			return
		} else {
			cmd = exec.Command(args[0], args[1:]...)
		}
		cmds = append(cmds, cmd)
	}

	if len(cmds) == 0 {
		return
	}

	for i := 0; i < len(cmds)-1; i++ {
		pipeReader, pipeWriter := io.Pipe()
		cmds[i].Stdout = pipeWriter
		cmds[i+1].Stdin = pipeReader
	}

	cmds[0].Stdin = os.Stdin
	cmds[len(cmds)-1].Stdout = os.Stdout

	for _, cmd := range cmds {
		cmd.Stderr = os.Stderr
	}

	for _, cmd := range cmds {
		err = cmd.Start()
		if err != nil {
			fmt.Fprintf(s.writer, "Error starting command: %v\n", err)
			s.lastError = err
			return
		}
	}

	for _, cmd := range cmds {
		err = cmd.Wait()
		if err != nil {
			s.lastError = err
			break
		}
	}
}

func (s *Shell) parseArguments(line string) []string {
	var args []string
	var current strings.Builder
	inQuotes := false
	quoteChar := byte(' ')

	for i := 0; i < len(line); i++ {
		c := line[i]

		switch {
		case c == '"' || c == '\'':
			if !inQuotes {
				inQuotes = true
				quoteChar = c
			} else if c == quoteChar {
				inQuotes = false
				if current.Len() > 0 {
					args = append(args, current.String())
					current.Reset()
				}
			} else {
				current.WriteByte(c)
			}

		case c == ' ' && !inQuotes:
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}

		default:
			current.WriteByte(c)
		}
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args
}

func main() {
	shell := NewShell()
	shell.Run()
}
