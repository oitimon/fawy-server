package server

type Config struct {
	Network     string `required:"true"`
	Host        string `required:"true"`
	Port        string `required:"true"`
	Timeout     int    `default:"10"`    // Timeout in seconds, 10 is default.
	MaxHandlers int    `default:"1000"`  // Pipeline size for handlers count.
	Difficulty  uint   `required:"true"` // PoW difficulty, must be from 1 to 100.
	Challenge   string `required:"true"` // PoW challenge driver.
}
