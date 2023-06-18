package enums

type calendarConsumer uint

const (
	CREATE calendarConsumer = iota
	UPDATE
	DELETE
	ADD_USER
	REMOVE_USER
	FETCH_CALENDAR
	AUTH
)

func (s calendarConsumer) String(id string) string {
	endpoint := "http://localhost:3021/calendar/v1"
	switch s {
	case CREATE:
		return endpoint + "/events"
	case FETCH_CALENDAR:
		return endpoint + "/events/" + id
	case UPDATE:
		return endpoint + "/events/" + id
	case DELETE:
		return endpoint + "/events/" + id
	case ADD_USER:
		return endpoint + "/events/" + id + "/users"
	case REMOVE_USER:
		return endpoint + "/events/" + id + "/users"
	default:
		return endpoint + "/auth"
	}
}
