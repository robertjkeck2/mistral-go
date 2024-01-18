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

func TestListModels(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client := mistral.NewMistralClient(os.Getenv("MISTRAL_API_KEY"))
	models, err := client.ListModels(context.Background())
	assert.Nil(t, err)
	assert.NotNil(t, models)
	assert.Contains(t, models.Data[0].ID, "mistral-")
}
