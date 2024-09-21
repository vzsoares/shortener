package types

type UrlBase struct {
	Rash        string `dynamodbav:"rash" json:"rash"`
	Destination string `dynamodbav:"destination" json:"destination"`
	Ttl         int    `dynamodbav:"ttl" json:"ttl"`
}

type UrlFull struct {
	*UrlBase
	UpdatedAt int `dynamodbav:"updatedAt" json:"updatedAt"`
	CreatedAt int `dynamodbav:"createdAt,unixtime" json:"createdAt"`
	Version   int `dynamodbav:"version" json:"version"`
}
