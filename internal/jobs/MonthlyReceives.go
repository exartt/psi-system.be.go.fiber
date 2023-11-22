package jobs

import (
	"github.com/robfig/cron/v3"
	"log"
	"psi-system.be.go.fiber/internal/repositories"
)

type MonthlyReceivesJob struct {
	repo repositories.TransactionRepository
}

func NewMonthlyReceivesJob(repo repositories.TransactionRepository) *MonthlyReceivesJob {
	return &MonthlyReceivesJob{
		repo: repo,
	}
}
func (j *MonthlyReceivesJob) Run() {
	log.Println("Running Monthly Receives Job")

	if err := j.repo.ThrowMonthlyReceives(); err != nil {
		log.Printf("Error running Monthly Receives Job: %v", err)
	}
}

func ScheduleJob(repo repositories.TransactionRepository) {
	c := cron.New()

	// para mostrar na banca. remover depois
	c.AddFunc("@every 1m", func() {
		//time.Sleep(1 * time.Minute)

		job := NewMonthlyReceivesJob(repo)
		job.Run()
	})

	//c.AddFunc("0 0 5 * *", func() {
	//	job := NewMonthlyReceivesJob(repo)
	//	job.Run()
	//})

	c.Start()
}
