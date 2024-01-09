package client

type Config struct {
	Host        string `required:"true"`
	Port        string `required:"true"`
	MaxRequests int    `default:"100"`   // Pipeline size for requests count.
	Timeout     int    `default:"10"`    // Timeout in seconds, 10 is default.
	Challenge   string `required:"true"` // PoW challenge driver.
}
