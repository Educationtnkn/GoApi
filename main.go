package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/gin-gonic/gin"
)

var dataMutex sync.Mutex

type Data struct {
	Sender    string `json:"sender"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

const jsonFileName = "chat.json"

func writeData(data Data) error {
	dataMutex.Lock()
	defer dataMutex.Unlock()

	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(jsonFileName, file, 0644)
	if err != nil {
		return err
	}

	return nil
}

func handleGet(c *gin.Context) {

	file, err := ioutil.ReadFile(jsonFileName)
	if err != nil {
		fmt.Println("Error reading file:", err)
		c.JSON(500, gin.H{"error": "Error reading messages"})
		return
	}

	var data []Data
	err = json.Unmarshal(file, &data)
	if err != nil {
		return
	}

	c.JSON(200, data)
}

func handlePost(c *gin.Context) {
	var newMessage Data
	err := c.ShouldBindJSON(&newMessage)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	// Read existing messages
	file, err := ioutil.ReadFile(jsonFileName)
	if err != nil {
		fmt.Println("Error reading file:", err)
		c.JSON(500, gin.H{"error": "Error reading messages"})
		return
	}

	var data []Data
	err = json.Unmarshal(file, &data)
	if err != nil {
		return
	}

	// Check if the message already exists, you can customize this logic based on your requirements
	// for _, msg := range data {
	// 	if msg.Sender == newMessage.Sender && msg.Content == newMessage.Content && msg.Timestamp == newMessage.Timestamp {
	// 		c.JSON(400, gin.H{"error": "Message already exists"})
	// 		return
	// 	}
	// }

	// Append the new message
	data = append(data, newMessage)

	filenew, err1 := json.MarshalIndent(data, "", "  ")
	if err1 != nil {
		return
	}

	err1 = ioutil.WriteFile(jsonFileName, filenew, 0644)
	if err1 != nil {
		return
	}

	c.Status(204)

}

func main() {
	r := gin.Default()
	r.Use(CORSMiddleware())
	api := r.Group("/api")
	{
		api.GET("/get", handleGet)
		api.POST("/post", handlePost)
	}

	fmt.Println("Server is running on http://localhost:8080")
	err := r.Run()
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
