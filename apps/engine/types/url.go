package types

type Url struct {
	Rash        string `dynamodbav:"rash" json:"rash"`
	Destination string `dynamodbav:"destination" json:"destination"`
	Ttl         int    `dynamodbav:"ttl" json:"ttl"`
	UpdatedAt   string `dynamodbav:"updatedAt" json:"updatedAt"`
	CreatedAt   string `dynamodbav:"createdAt,unixtime" json:"createdAt"`
}
