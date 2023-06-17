package repositories

import (
	"strconv"
	"time"

	"gorm.io/gorm"
	"psi-system.be.go.fiber/internal/domain/enums"
	"psi-system.be.go.fiber/internal/domain/model/appointment"
)

type AppointmentRepository interface {
	Create(appointment *appointment.Appointment) error
	Delete(ID uint) error
	Update(ID uint, appointment *appointment.Appointment) error
	GetAll() ([]*appointment.Appointment, error)
	GetByID(ID uint) (*appointment.Appointment, error)
	GetByStart(start time.Time) (*appointment.Appointment, error)
	GetByTimeRange(start, end time.Time) ([]*appointment.Appointment, error)
	GetByTimeRangeNotId(id uint, start, end time.Time) ([]*appointment.Appointment, error)
	CheckDateByTimeRange(start, end time.Time) ([]*appointment.Appointment, error)
	CheckDateByTimeRangeUpdate(start, end time.Time) ([]*appointment.Appointment, error)
	UpdateStatusAppointment(ID uint, status enums.StatusAgendamento) error
	GetByYear(year int) ([]*appointment.Appointment, error)
}

type appointmentRepository struct {
	db *gorm.DB
}

func NewAppointmentRepository(db *gorm.DB) AppointmentRepository {
	return &appointmentRepository{db: db}
}

func (r *appointmentRepository) GetByYear(year int) ([]*appointment.Appointment, error) {
	var appointments []*appointment.Appointment
	startOfYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	endOfYear := time.Date(year, 12, 31, 23, 59, 59, 999999999, time.UTC)

	canceladoStatus := strconv.Itoa(int(enums.Cancelado))

	err := r.db.Where("\"start\" >= ? AND \"start\" <= ? AND \"status\" != ?", startOfYear, endOfYear, canceladoStatus).Find(&appointments).Error
	return appointments, err
}

func (r *appointmentRepository) Create(appointment *appointment.Appointment) error {
	result := r.db.Create(appointment)
	return result.Error
}

func (r *appointmentRepository) Delete(ID uint) error {
	result := r.db.Delete(&appointment.Appointment{}, ID)
	return result.Error
}

func (r *appointmentRepository) Update(ID uint, appointment *appointment.Appointment) error {
	result := r.db.Model(appointment).Where("id = ?", ID).Save(appointment)
	return result.Error
}

func (r *appointmentRepository) UpdateStatusAppointment(ID uint, status enums.StatusAgendamento) error {
	result := r.db.Model(&appointment.Appointment{}).Where("id = ?", ID).Update("status", status)
	return result.Error
}

func (r *appointmentRepository) GetAll() ([]*appointment.Appointment, error) {
	var appointments []*appointment.Appointment
	err := r.db.Find(&appointments).Error
	return appointments, err
}

func (r *appointmentRepository) GetByID(ID uint) (*appointment.Appointment, error) {
	var appointment appointment.Appointment
	err := r.db.First(&appointment, "id = ?", ID).Error
	return &appointment, err
}

func (r *appointmentRepository) GetByStart(start time.Time) (*appointment.Appointment, error) {
	var appointment appointment.Appointment
	err := r.db.First(&appointment, "start = ?", start).Error
	return &appointment, err
}

func (r *appointmentRepository) GetByTimeRange(start, end time.Time) ([]*appointment.Appointment, error) {
	var appointments []*appointment.Appointment
	canceladoStatus := strconv.Itoa(int(enums.Cancelado))
	err := r.db.Not("status = ?", canceladoStatus).Where("\"start\" >= ? AND \"end\" < ?", start, end).Find(&appointments).Error
	return appointments, err
}

func (r *appointmentRepository) GetByTimeRangeNotId(ID uint, start, end time.Time) ([]*appointment.Appointment, error) {
	var appointments []*appointment.Appointment
	canceladoStatus := int(enums.Cancelado)
	err := r.db.Not("id = ?", ID).Not("status = ?", canceladoStatus).Where("\"start\" >= ? AND \"end\" < ?", start, end).Find(&appointments).Error
	return appointments, err
}

func (r *appointmentRepository) CheckDateByTimeRange(start, end time.Time) ([]*appointment.Appointment, error) {
	var appointments []*appointment.Appointment
	err := r.db.Where("\"start\" < ? AND \"end\" > ?", end, start).Find(&appointments).Error
	return appointments, err
}

func (r *appointmentRepository) CheckDateByTimeRangeUpdate(start, end time.Time) ([]*appointment.Appointment, error) {
	var appointments []*appointment.Appointment
	err := r.db.Where("\"start\" < ? AND \"end\" > ?", end, start).Find(&appointments).Error
	return appointments, err
}
