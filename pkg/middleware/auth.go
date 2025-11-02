package middleware

import(
	
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header=="" {
			c.AbortWithStatusJSON(http.StatusUnanthorized,gin.H{"error":"missing token"})
			return
		}

		parts := strings.SplitN(h,"",2)
		if !(len(parts)==2 && parts[0]=="Bearer"){
			c.AbortWithStatusJSON(http.StatusUnanthorized,gin.H{"error":"invalid auto header"})
			return
		}

		claims,err := jwt.ParseToken(parts[1])
		if err!=nil {
			switch err {
				case 
			}
		}
	}
}