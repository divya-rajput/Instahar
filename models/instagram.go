package models

import "encoding/json"

type InstagramMedia struct {
	TakenAt         int64           `json:"taken_at"`
	DeviceTimestamp int64           `json:"device_timestamp"`
	MediaType       int             `json:"media_type"`
	UserTags        json.RawMessage `json:"user_tags"`
	LikeCount       int64           `json:"like_count"`
	Code            string          `json:"code"`
	Caption         struct {
		Type int
		Text string
	} `json:"caption"`
	PlayCount         int64           `json:"play_count"`
	CoauthorProducers json.RawMessage `json:"coauthor_producers"`
	ProductType       string          `json:"product_type"`
	Comments          json.RawMessage `json:"comments"`
	CommentCount      int64           `json:"comment_count"`
	ClipsMetadata     struct {
		MusicInfo         json.RawMessage `json:"music_info"`
		OriginalSoundInfo json.RawMessage `json:"original_sound_info"`
		AudioType         json.RawMessage `json:"audio_type"`
		MashupInfo        json.RawMessage `json:"mashup_info"`
	} `json:"clips_metadata"`
}

type InstagramMediaEntry struct {
	Media InstagramMedia `json:"media"`
}

type InstagramGraphQLQueryResponse struct {
	Items []InstagramMediaEntry `json:"items"`
}
