package person

import (
	"encoding/json"
	"fmt"
)

type PersonDTO struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

func UnserializePersonDto(data []byte) (*PersonDTO, error) {
	var res PersonDTO
	err := json.Unmarshal(data, &res)

	if err != nil {
		return nil, fmt.Errorf("failed to unserialize person: %s", err.Error())
	}

	return &res, nil
}

func (m *PersonDTO) Serialized() ([]byte, error) {
	res, err := json.Marshal(*m)

	if err != nil {
		return nil, fmt.Errorf("failed to serialize person: %v", err)
	}

	return res, nil
}
