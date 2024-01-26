package mistral_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/robertjkeck2/mistral-go"
	"github.com/stretchr/testify/assert"
)

func TestCreateChatCompletion(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client := mistral.NewMistralClient(os.Getenv("MISTRAL_API_KEY"))
	completion, err := client.CreateChatCompletion(context.Background(), mistral.ChatCompletionRequest{
		Model:    "mistral-tiny",
		Messages: []mistral.ChatMessage{{Role: mistral.RoleUser, Content: "Respond to this message with the message: \"Hello, world!\""}},
	})
	assert.Nil(t, err)
	assert.NotNil(t, completion)
	assert.Equal(t, "chat.completion", completion.Object)
	assert.Equal(t, "Hello, world!", completion.Choices[0].Message.Content)
}
