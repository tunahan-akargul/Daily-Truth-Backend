package words

import "time"

type Word struct {
	ID        string    `bson:"_id.omitempty" json:"id"`
	Text      string    `bson:"text" json:"text"`
	OwnerID   string    `bson:"ownerid" json:"ownerId"`
	CreatedAt time.Time `bson:"createdat" json:"createdAt"`
}

type CreateWordReq struct {
	Text string `json:"text"`
}
