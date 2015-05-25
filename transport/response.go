package transport

/* The serialized format:
 *    status:uint8
 */

type Response struct {
	Status  ResponseStatus
	Payload []byte
}

func Unmarshal(buf []byte) *Response {
	status, payload := buf[1], buf[1:]
	return &Response{
		Status:  ResponseStatus(status),
		Payload: payload,
	}
}
