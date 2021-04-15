package helper

import "encoding/json"

type response struct {
	Message string `json:"message"`
}

func ResponseJSON(message string) ([]byte, error) {
	responseStruct := response{Message: message}
	responseJson, erro := json.Marshal(responseStruct)

	if erro != nil {
		return nil, erro
	}

	return responseJson, nil
}

func RespMessageError(message string) []byte {
	responseStruct := response{Message: message}
	responseJson, _ := json.Marshal(responseStruct)
	return responseJson
}
