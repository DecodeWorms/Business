package types

import "time"

type TokenDetails struct {
	AccessToken  string
	AtExp        time.Time
	RefreshToken string
	RfExp        time.Time
}
