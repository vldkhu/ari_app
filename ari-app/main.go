package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"
)

var (
	ctx         = context.Background()
	redisClient *redis.Client
	ariURL      = "http://127.0.0.1:8088/ari"
	ariUser     = "admin"
	ariPass     = "admin"
	channel     = "SIP/sip_account"
)

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	// Инициализация данных
	soundFiles := []string{"sound1.wav", "sound2.wav", "sound3.wav"}
	for _, file := range soundFiles {
		err := redisClient.LPush(ctx, "sound_files", file).Err()
		if err != nil {
			log.Fatalf("Ошибка при добавлении файла %s: %v", file, err)
		}
	}

	fmt.Println("Данные успешно инициализированы в Redis.")
}

// getRandomSoundFile получает случайный звуковой файл из Redis и удаляет его из списка
func getRandomSoundFile() (string, error) {
	soundFiles, err := redisClient.LRange(ctx, "sound_files", 0, -1).Result()
	if err != nil {
		return "", err
	}

	if len(soundFiles) == 0 {
		return "", nil
	}

	randomIndex := rand.Intn(len(soundFiles))
	soundFile := soundFiles[randomIndex]

	// Удаляем файл из Redis
	redisClient.LRem(ctx, "sound_files", 1, soundFile)

	log.Printf("Получен звуковой файл: %s", soundFile)
	return soundFile, nil
}

// отправляет запрос на воспроизведение звукового файла через ARI
func playSoundFile(soundFile string) error {
	client := resty.New()
	_, err := client.R().
		SetBasicAuth(ariUser, ariPass).
		SetBody(map[string]string{"media": "sound:" + soundFile}).
		Post(fmt.Sprintf("%s/channels/%s/play", ariURL, channel))
	if err != nil {
		log.Printf("Ошибка при воспроизведении звукового файла %s: %v", soundFile, err)
	} else {
		log.Printf("Воспроизведение звукового файла %s началось", soundFile)
	}
	return err
}

func hangupCall() {
	client := resty.New()
	_, err := client.R().
		SetBasicAuth(ariUser, ariPass).
		Post(fmt.Sprintf("%s/channels/%s/hangup", ariURL, channel))
	if err != nil {
		log.Printf("Ошибка при сбросе звонка: %v", err)
	} else {
		log.Println("Звонок сброшен")
	}
}

// обрабатывает входящие звонки, получает случайный звуковой файл и воспроизводит его
func handleIncomingCall(w http.ResponseWriter, r *http.Request) {
	soundFile, err := getRandomSoundFile()
	if err != nil {
		log.Printf("Ошибка при получении звукового файла: %v", err)
		hangupCall()
		return
	}

	if soundFile == "" {
		log.Println("Нет доступных звуковых файлов, сбрасываем звонок.")
		hangupCall()
		return
	}

	err = playSoundFile(soundFile)
	if err != nil {
		log.Printf("Ошибка при воспроизведении звукового файла: %v", err)
		hangupCall()
	}

	time.AfterFunc(5*time.Second, func() {
		hangupCall()
	})
}

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	http.HandleFunc("/incoming", handleIncomingCall)
	fmt.Println("Starting ARI application...")
	if err := http.ListenAndServe(":8088", nil); err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
}
