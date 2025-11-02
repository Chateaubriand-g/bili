package config

type Config struct {
	JwtSecret string
}

func LoadConfig() *Config {
	//godotenv会捕捉工作目录下的.env文件，将文件中的变量设置到系统变量中
	if err:= godotenv.Load();err!=nil{
		log.Println("No .env file found,using system environment variales")
	}

	return &Config{
		JwtSecret: getenv("JWT_SECRET")
	}
}

func getenv(key,value string) string {
	//查询对应系统变量的value及exists
	if value,exists := os.LookupEnv(key);!exists{
		log.Fatal("FATAL ERROR:%s not set",value)
	}
	return value
}