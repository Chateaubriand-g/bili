package jwt

var (
	expiredHour = 2 
	waitTime = 0
) 

var (
	ErrTokenExpired = errors.New("token expired")
	ErrTokenInvalid = errors.New("token invaild")
	ErrTokenNotVaildYet = errors.New("token not active yet")
	ErrSignatureInvalid = errors.New("signatrue invaild")
)

type Claims struct {
	UserID uint64 `json:"userId"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint64,secret string) (string,error) {
	expiredTime = time.Now().Add(expiredHour*time.Hour)

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredTime),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()+waitTime*time.Hour),
			Issuer: "g",
			Subject: fmt.Sprintf("%d",userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	tokenstring,err := jwt.SignedString([]byte(secret))
	if err != nil {
		return "",err
	}
	return tokenstring,nil
}

func ParseToken(tokenstring,secret string) (*Claims,error) {
	/*
	jwt.Token{
		Raw string
		Method SigningMethod
		Header map[string]interface{}
		Claim Claims
		Signature string
		Valid bool
	}	
	*/
	token,err := jwt.ParseWithClaims(tokenstring,&Claims{},func(token *jwt.Token) (interface{},error){
		if _,ok :=token.Method.(*jwt.SigningMethodHMAC);!ok {
			return nil,fmt.Errorf("unexpected signing method: %v",token.Header["alg"])
		}
		return []byte(secret),nil
	})

	if err!=nil {
		if errors.Is(err,jwt.ErrTokenExpired){
			return nil,ErrTokenExpired
		}else if errors.Is(err,jwt.ErrTokenNotValidYet){
			return nil,ErrTokenNotValidYet
		}else if errors.Is(err,jwt.ErrSignatureInvalid){
			return nil,ErrSignatureInvalid
		}
		return nil,ErrTokenInvalid
	}

	if claims,ok := token.Claims.(*Claims);ok && token.Valid {
		return claims,nil
	}
	return nil,ErrTokenInvalid
}