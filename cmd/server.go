package main

func main(){
	router :=gin.Default()
	router.GET("/", func(ctx *gin.Context){
		ctx.JSON(200, gin.H{"msg": "hello world~~~"})
	})
	_ = router.Run("0.0.0.0:8000")
}
