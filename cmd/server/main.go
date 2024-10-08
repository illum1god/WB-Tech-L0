package main

import (
	kafkaGateway "WB-L0/internal/gateways/kafka"
	"context"
	"errors"
	_ "github.com/lib/pq"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"
	"time"

	"WB-L0/internal/configs"
	httpGateway "WB-L0/internal/gateways/http"
	"WB-L0/internal/repository"
	"WB-L0/internal/repository/order/postgres"
	"WB-L0/internal/usecase"

	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	absPath, err := filepath.Abs(".")
	if err != nil {
		logger.Fatalf("Не удалось получить абсолютный путь: %v", err)
	}
	logger.Infof("Абсолютный путь: %s", absPath)

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "internal/configs/config.yaml"
	}

	config, err := configs.LoadConfig(configPath)
	if err != nil {
		logger.Fatalf("Не удалось загрузить конфигурацию: %v", err)
	}

	logger.Infof("Подключение к базе данных на %s:%s", config.DB.Host, config.DB.Port)

	db, err := postgres.NewPostgresDB(configs.DBConfig{
		Host:     config.DB.Host,
		Port:     config.DB.Port,
		Username: config.DB.Username,
		DBName:   config.DB.DBName,
		SSLMode:  config.DB.SSLMode,
		Password: config.DB.Password,
	})
	if err != nil {
		logger.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Errorf("Ошибка при закрытии базы данных: %v", err)
		}
	}()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	repos := repository.NewRepository(db)

	err = repos.RestoreCache()
	if err != nil {
		log.Fatalf("Не удалось восстановить кэш из базы данных: %v", err)
	}

	services := usecase.NewService(repos)

	time.Sleep(20 * time.Second)
	kafkaClient := kafkaGateway.NewConsumer(config.Kafka, services)

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	err = waitForKafka(config.Kafka.Brokers, 30*time.Second)
	if err != nil {
		logger.Fatalf("Kafka недоступна: %v", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := kafkaClient.Consume(&wg, ctx); err != nil {
			logger.Fatalf("Не удалось подписаться на топик: %v", err)
		}
	}()
	logger.Infof("Подписка на Kafka topic: %s", config.Kafka.Topic)

	host := config.Server.Host
	port, err := strconv.Atoi(config.Server.Port)
	if err != nil {
		logger.Warnf("Некорректный порт сервера, используется 8080 по умолчанию: %v", err)
		port = 8080
	}

	server := httpGateway.NewServer(
		services,
		httpGateway.WithHost(host),
		httpGateway.WithPort(uint16(port)),
	)

	go func() {
		logger.Infof("HTTP сервер запущен на %s:%d", host, port)
		if err := server.Run(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				logger.Fatalf("Ошибка запуска HTTP сервера: %v", err)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Остановка сервера...")

	cancel()
	wg.Wait()

	kafkaClient.Stop()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Errorf("Ошибка при остановке HTTP сервера: %v", err)
	}

	logger.Info("Сервер успешно остановлен")
}

func waitForKafka(brokers []string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		for _, broker := range brokers {
			conn, err := net.DialTimeout("tcp", broker, 2*time.Second)
			if err != nil {
				log.Printf("Не удалось подключиться к брокеру Kafka %s: %v", broker, err)
				break
			}
			_ = conn.Close()
		}
		return nil
	}
	return errors.New("kafka недоступна после ожидания")
}
