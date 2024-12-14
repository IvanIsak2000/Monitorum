package utils

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
)

type Config struct{
	Token string
	ModerId telego.ChatID
}
func GetConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Ошибка загрузки .env файла:", err)
	}
	moderID := os.Getenv("MODER_ID")
	moderIDInt, err := strconv.ParseInt(moderID, 10, 64)
	if err != nil {
		log.Fatalf("Ошибка преобразования MODER_ID: %v\n", err)
		return Config{}, fmt.Errorf("Ошибка преобразования MODER_ID: %v\n", err)
	}
	return Config{
		Token: os.Getenv("TOKEN"),
		ModerId: telego.ChatID{ID: moderIDInt},
	}, nil

}