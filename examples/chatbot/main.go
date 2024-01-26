package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/robertjkeck2/mistral-go"
)

var (
	AvailableModels   = []string{"mistral-tiny", "mistral-small", "mistral-medium"}
	AvailableCommands = []string{"/new", "/model", "/system", "/temperature", "/safemode", "/config", "/exit", "/quit", "/help"}
)

// ChatBot is an interactive chatbot that uses the Mistral Platform to generate messages
type ChatBot struct {
	Client    mistral.MistralClient
	Model     string
	SystemMsg mistral.ChatMessage
	Temp      float64
	SafeMode  bool
	Messages  []mistral.ChatMessage
}

func (c *ChatBot) Start() {
	c.ShowInstructions()
	c.New()
	c.Run()
}

func (c *ChatBot) Exit() {
	fmt.Println("Exiting chat...")
	os.Exit(0)
}

func (c *ChatBot) ShowInstructions() {
	fmt.Println(`To chat: type your message and hit enter
To start a new chat: /new
To switch model: /model <model name>
To switch system message: /system <message>
To switch temperature: /temperature <temperature>
To switch safe mode: /safemode <true/false>
To see current config: /config
To exit: /exit, /quit, or hit CTRL+C
To see this help: /help`)
}

func (c *ChatBot) New() {
	fmt.Println("")
	fmt.Println("Starting new chat with model: " + c.Model + ", temperature: " + fmt.Sprintf("%f", c.Temp) + ", safe mode: " + fmt.Sprintf("%t", c.SafeMode))
	fmt.Println("")

	c.Messages = []mistral.ChatMessage{}
	if c.SystemMsg.Content != "" {
		c.Messages = append(c.Messages, c.SystemMsg)
	}
}

func (c *ChatBot) UpdateModel(model string) {
	if model == c.Model {
		fmt.Println("Model is already: " + c.Model)
		return
	}
	if !slices.Contains(AvailableModels, model) {
		fmt.Println("Invalid model: " + model)
		return
	}
	c.Model = model
	fmt.Println("Updated model to: " + c.Model)
}

func (c *ChatBot) UpdateTemp(tempStr string) {
	temp, err := strconv.ParseFloat(tempStr, 64)
	if err != nil {
		fmt.Println("Invalid temperature: " + tempStr)
		return
	}
	if temp < 0.0 || temp > 1.0 {
		fmt.Println("Invalid temperature: " + fmt.Sprintf("%f", temp))
		return
	}
	c.Temp = temp
	fmt.Println("Updated temperature to: " + fmt.Sprintf("%f", c.Temp))
}

func (c *ChatBot) UpdateSafeMode(safeMode string) {
	safeModeBool, err := strconv.ParseBool(safeMode)
	if err != nil {
		fmt.Println("Invalid safe mode: " + safeMode)
		return
	}
	if safeModeBool == c.SafeMode {
		fmt.Println("Safe mode is already: " + fmt.Sprintf("%t", c.SafeMode))
		return
	}
	c.SafeMode = safeModeBool
	fmt.Println("Updated safe mode to: " + fmt.Sprintf("%t", c.SafeMode))
}

func (c *ChatBot) UpdateSystemMsg(systemMsg string) {
	if systemMsg == c.SystemMsg.Content {
		fmt.Println("System message is already: " + c.SystemMsg.Content)
		return
	}
	c.SystemMsg.Content = systemMsg
	c.Messages = []mistral.ChatMessage{c.SystemMsg}
	fmt.Println("Updated system message to: " + c.SystemMsg.Content)
}

func (c *ChatBot) ShowConfig() {
	fmt.Println("")
	fmt.Println("Current model: " + c.Model)
	fmt.Println("Current temperature: " + fmt.Sprintf("%f", c.Temp))
	fmt.Println("Current safe mode: " + fmt.Sprintf("%t", c.SafeMode))
	fmt.Println("Current system message: " + c.SystemMsg.Content)
	fmt.Println("")
}

func (c *ChatBot) GetInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			return line
		}
	}
	return ""
}

func (c *ChatBot) SendMessage(message string) {
	c.Messages = append(c.Messages, mistral.ChatMessage{Role: mistral.RoleUser, Content: message})
	chatCompletionRequest := mistral.ChatCompletionRequest{
		Model:       c.Model,
		Messages:    c.Messages,
		Temperature: c.Temp,
		SafeMode:    c.SafeMode,
	}
	fmt.Println(c.Messages)
	chatCompletionStream, err := c.Client.CreateChatCompletionStream(context.Background(), chatCompletionRequest)
	if err != nil {
		panic(err)
	}

	defer chatCompletionStream.Close()

	assistantMessage := []string{}
	for {
		var response mistral.ChatCompletionStreamResponse
		response, err = chatCompletionStream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Printf(response.Choices[0].Delta.Content)
		assistantMessage = append(assistantMessage, response.Choices[0].Delta.Content)
	}

	c.Messages = append(c.Messages, mistral.ChatMessage{Role: mistral.RoleAssistant, Content: strings.Join(assistantMessage, "")})
	fmt.Println("")
	fmt.Println("")
}

func (c *ChatBot) Run() {
	for {
		command := c.GetInput()
		if command == "" {
			continue
		}
		commandSplit := strings.Split(command, " ")
		if slices.Contains(AvailableCommands, commandSplit[0]) {
			args := []string{}
			if len(commandSplit) > 1 {
				args = commandSplit[1:]
			}
			c.Execute(commandSplit[0], args)
		} else {
			c.SendMessage(command)
		}
	}
}

func (c *ChatBot) Execute(command string, args []string) {
	switch command {
	case "/new":
		c.New()
	case "/model":
		c.UpdateModel(args[0])
	case "/system":
		c.UpdateSystemMsg(strings.Join(args, " "))
	case "/temperature":
		c.UpdateTemp(args[0])
	case "/safemode":
		c.UpdateSafeMode(args[0])
	case "/exit", "/quit":
		c.Exit()
	case "/help":
		c.ShowInstructions()
	case "/config":
		c.ShowConfig()
	default:
		fmt.Println("Invalid command: " + command)
	}
}

func main() {
	model := flag.String("model", "mistral-tiny", "Mistral model to use for the chatbot")
	systemMsg := flag.String("systemMsg", "", "System message to send to the chatbot at the beginning of the chat")
	temp := flag.Float64("temp", 0.7, "Temperature to use for the chatbot")
	safeMode := flag.Bool("safeMode", false, "Whether to use safe mode for the chatbot")

	flag.Parse()

	chatbot := ChatBot{
		Client:    *mistral.NewMistralClient(os.Getenv("MISTRAL_API_KEY")),
		Model:     *model,
		SystemMsg: mistral.ChatMessage{Role: mistral.RoleSystem, Content: *systemMsg},
		Temp:      *temp,
		SafeMode:  *safeMode,
		Messages:  []mistral.ChatMessage{},
	}

	chatbot.Start()
}
