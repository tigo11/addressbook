
package stdhttp

import (
	"encoding/json"
	"log"
	"hw2/gate/psg"
	"hw2/models/dto"
	"net/http"
	"hw2/pkg"
)

type Controller struct {
	DB  *psg.Psg
	Srv *http.Server
}

func NewController(addr string, db *psg.Psg) *Controller {
	return &Controller{
		DB: db,
		Srv: &http.Server{Addr: addr},
	}
}

// обработка ошибок
func handleError(w http.ResponseWriter, errorMsg string, statusCode int, err error) {
	http.Error(w, errorMsg, statusCode)
	log.Printf("%s: %v", errorMsg, err)
}

// RecordAdd добавляет новую запись в адресную книгу.
func (c *Controller) RecordAdd(w http.ResponseWriter, r *http.Request) {
	var record dto.Record
	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		handleError(w, "Invalid request body", http.StatusBadRequest, err)
		return
	}

	normalizedPhone, err := pkg.PhoneNormalize(record.Phone)
	if err != nil {
		handleError(w, "Error normalizing phone number", http.StatusInternalServerError, err)
		return
	}
	record.Phone = normalizedPhone

	id, err := c.DB.RecordAdd(record)
	if err != nil {
		handleError(w, "Error adding record to the database", http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}

// RecordsGet извлекает записи из адресной книги на основе заданных параметров.
func (c *Controller) RecordsGet(w http.ResponseWriter, r *http.Request) {
	var record dto.Record
	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		handleError(w, "Ошибка декодирования запроса", http.StatusBadRequest, err)
		return
	}

	records, err := c.DB.RecordsGet(record)
	if err != nil {
		handleError(w, "Ошибка выполнения запроса к базе данных", http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(records); err != nil {
		handleError(w, "Ошибка кодирования ответа", http.StatusInternalServerError, err)
	}
}

// RecordUpdate обновляет запись в адресной книге.
func (c *Controller) RecordUpdate(w http.ResponseWriter, r *http.Request) {
	var record dto.Record
	err := json.NewDecoder(r.Body).Decode(&record)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	normalizedPhone, err := pkg.PhoneNormalize(record.Phone)
	if err != nil {
		http.Error(w, "Error normalizing phone number", http.StatusInternalServerError)
		return
	}
	record.Phone = normalizedPhone

	err = c.DB.RecordUpdate(record)
	if err != nil {
		http.Error(w, "Error updating record in the database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// RecordDeleteByPhone удаляет запись по номеру телефона.
func (c *Controller) RecordDeleteByPhone(w http.ResponseWriter, r *http.Request) {
	var record dto.Record
	err := json.NewDecoder(r.Body).Decode(&record)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	normalizedPhone, err := pkg.PhoneNormalize(record.Phone)
	if err != nil {
		http.Error(w, "Error normalizing phone number", http.StatusInternalServerError)
		return
	}
	record.Phone = normalizedPhone

	err = c.DB.RecordDeleteByPhone(record.Phone)
	if err != nil {
		http.Error(w, "Error deleting record from the database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
