package mistral_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/robertjkeck2/go-mistral"
	"github.com/stretchr/testify/assert"
)

func TestCreateEmbedding(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client := mistral.NewMistralClient(os.Getenv("MISTRAL_API_KEY"))
	embeddings, err := client.CreateEmbedding(context.Background(), mistral.EmbeddingRequest{
		Model: "mistral-embed",
		Input: []string{"hello", "world"},
	})
	assert.Nil(t, err)
	assert.NotNil(t, embeddings)
	assert.Equal(t, "embedding", embeddings.Data[0].Object)
}
