package cron

import (
	"bytes"
	"log"
	"log/slog"

	"go-tg-bot/internal/repository"
	"go-tg-bot/internal/services"
	"go-tg-bot/internal/transport"

	"github.com/robfig/cron/v3"
	"gopkg.in/telebot.v4"
)

type Scheduler struct {
	cron         *cron.Cron
	services     *services.Services
	repositories *repository.Repositories
	log          *slog.Logger
	tg           *transport.Telegram
}

func NewScheduler(services *services.Services, repositories *repository.Repositories, tg *transport.Telegram, log *slog.Logger) *Scheduler {
	scheduler := &Scheduler{
		cron:         cron.New(cron.WithSeconds()),
		services:     services,
		repositories: repositories,
		log:          log,
		tg:           tg,
	}
	scheduler.Start()

	return scheduler
}

func (s *Scheduler) Start() {
	_, err := s.cron.AddFunc("0 00 23 * * *", func() {
		messages := s.repositories.Message.GetMessagesForToday()
		if len(messages) == 0 {
			log.Println("Сегодня сообщений нет")
			return
		}
		content := s.services.Messages.FormatMessagesForGigaChat(messages)

		result := s.services.GigaChat.Send(content)

		txt := result.Choices[0].Message.Content

		s.repositories.Report.Create(txt)

		imgData, e := s.services.GigaChat.GenerateImage("Нарисуй изображение по отчету из нашено чата: " + txt)

		bot := s.tg.Bot()
		if e != nil || len(imgData) == 0 {
			// Если нет изображения — отправляем только текст
			_, e := bot.Send(telebot.ChatID(-4204971428), txt)
			if e != nil {
				log.Println("Ошибка при отправке текста:", e)
			}
			log.Println("DONE (только текст)")
			return
		}

		photo := &telebot.Photo{
			File:    telebot.FromReader(bytes.NewReader(imgData)),
			Caption: txt,
		}
		log.Println("DONE!!!!!!!!!")
		_, er := bot.Send(telebot.ChatID(-4204971428), photo)
		if er != nil {
			log.Fatal("Ошибка получения последнего отчета:", er)
		}
	})
	if err != nil {
		s.log.Error("ошибка при добавлении задачи cron", "err", err)
	}

	s.cron.Start()
}

func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
}
