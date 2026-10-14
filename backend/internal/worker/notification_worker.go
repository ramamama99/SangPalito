package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// NotificationType represents the event triggering the notification.
type NotificationType string

const (
	NotifyAppointmentCreated NotificationType = "appointment_created"
	NotifyQueueCalled        NotificationType = "queue_called"
	NotifyQueueApproaching   NotificationType = "queue_approaching"
)

// NotificationJob contains payload details for sending SMS/WhatsApp messages.
type NotificationJob struct {
	Type        NotificationType
	QueueID     uuid.UUID
	QueueNumber int
	Recipient   string // Phone number
	PatientName string
	VisitDate   time.Time
}

// NotificationWorker dispatches asynchronous SMS and WhatsApp notifications.
type NotificationWorker struct {
	jobs chan NotificationJob
}

// NewNotificationWorker creates an asynchronous notification worker.
func NewNotificationWorker(bufferSize int) *NotificationWorker {
	if bufferSize <= 0 {
		bufferSize = 200
	}
	return &NotificationWorker{
		jobs: make(chan NotificationJob, bufferSize),
	}
}

// Enqueue puts a notification job into the channel without blocking caller.
func (w *NotificationWorker) Enqueue(job NotificationJob) {
	if w == nil {
		return
	}
	select {
	case w.jobs <- job:
	default:
		slog.Warn("notification worker buffer full, dropping job", "type", job.Type, "queue_id", job.QueueID)
	}
}

// Start listens for notification jobs and processes them in the background.
func (w *NotificationWorker) Start(ctx context.Context) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("notification worker panic recovered", "panic", r)
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case job := <-w.jobs:
				w.processJob(job)
			}
		}
	}()
}

func (w *NotificationWorker) processJob(job NotificationJob) {
	slog.Info("notification dispatched async",
		"type", job.Type,
		"recipient", job.Recipient,
		"patient_name", job.PatientName,
		"queue_number", job.QueueNumber,
	)
}
