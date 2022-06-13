package messaging

// EventPersonalData represents the `personaldata` field
type EventPersonalData struct {
	Initiator          EventPersonalDataInitiator          `json:"initiator"`
	AnonymizeURL       EventPersonalDataAnonymizeURL       `json:"anonymizeurl"`
	RelatedDataEntries []EventPersonalDataRelatedDataEntry `json:"relateddataentries"`
}

// EventPersonalDataInitiator represents the `personaldata`->`initiator` field
type EventPersonalDataInitiator struct {
	AccessTokenID     string `json:"access-token-id"`
	AccessTokenName   string `json:"access-token-name"`
	AccessTokenEmail  string `json:"access-token-email"`
	StoreID           string `json:"store-id"`
	CustomerTokenID   string `json:"customer-token-id"`
	CustomerTokenName string `json:"customer-token-name"`
}

// EventPersonalDataAnonymizeURL represents the `personaldata`->`anonymizeurl` field
type EventPersonalDataAnonymizeURL struct {
	Method string      `json:"method"`
	URL    string      `json:"url"`
	Body   interface{} `json:"body"`
}

// EventPersonalDataRelatedDataEntry represents the `personaldata`->`relateddataentries` field
type EventPersonalDataRelatedDataEntry struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}
