package pkg

import (
	"encoding/json"
)

type User struct {
	Id    string `json:"id"`
	Token string `json:"token"`
	//Queue         int    `json:"queue"`           // Henüz yanıt döndürmemiş istekler için
	CanNotRequest bool `json:"can_not_request"` // Başka bir processes tarafından denetlenmeli diye düşünüyorum

	// Diğer değişkenler proxy'ın pek ilgisini çekmiyor. Fazla veri ile deserialazyon yapmaya gerek yok.
}

func (u *User) UnmarshalJSON(data []byte) error {
	// fastjson.ParseBytes(data) // Çekilecek veri çok az, gereksiz kaçar
	return json.Unmarshal(data, u)

}

func (u User) ToMapI() (m map[string]interface{}) {
	return map[string]interface{}{
		"id":              u.Id,
		"token":           u.Token,
		"can_not_request": u.CanNotRequest,
	}

}

func UserFromMapI(m map[string]interface{}) (u User) {
	//fmt.Printf("m: %+v\n", m)
	u.Id, _ = m["id"].(string)
	u.Token, _ = m["token"].(string)
	u.CanNotRequest, _ = m["can_not_request"].(bool)

	return

}
