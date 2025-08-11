package services

type VoiceOfReason struct {
	Reason   string
	UserName string
	Position string
	NumLikes int
}
type WagerDecisionDetails struct {
	Title          string
	Description    string
	Left           string
	Right          string
	VoicesOfReason []VoiceOfReason
}

type WagerVerdict struct {
}

func Decide(wagerDecisionDetails WagerDecisionDetails) (string, string, error) {
	return "", "", nil
}
