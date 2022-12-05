package pkg

import (
	"context"
	"encoding/json"
)

type User struct {
	UserDatabase `json:"-"`

	Id    string `json:"id"`
	Token string `json:"token"`
	//Queue         int    `json:"queue"`           // Henüz yanıt döndürmemiş istekler için
	CanNotRequest bool `json:"can_not_request"` // Başka bir processes tarafından denetlenmeli diye düşünüyorum

	// Diğer değişkenler proxy'ın pek ilgisini çekmiyor. Fazla veri ile deserialazyon yapmaya gerek yok.
}

func (u User) Update(ctx context.Context) error {
	return u.UserDatabase.UpdateUser(ctx, u)
}

func (u *User) UnmarshalJSON(data []byte) error {
	// fastjson.ParseBytes(data) // Çekilecek veri çok az, gereksiz kaçar
	return json.Unmarshal(data, u)

}
