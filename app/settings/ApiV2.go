package settings

var ApiV2 = &apiV2{
	ClientId:     5, // osu!lazer
	ClientSecret: "FGc9GAtyHzeQDshWP5Ah7dega8hJACAJpQtw6OXk",
	Username:     "Guest",
	Password:     "",
}

type apiV2 struct {
	ClientId     int
	ClientSecret string
	Username     string
	Password     string `password:"true"`
}
