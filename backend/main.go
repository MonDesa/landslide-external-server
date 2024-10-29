package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	configDir := "./configs"

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		os.Mkdir(configDir, os.ModePerm)
	}

	router.GET("/getConfig", func(c *gin.Context) {
		commUnitID := c.Query("CommUnitID")
		if commUnitID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "CommUnitID not provided"})
			return
		}

		configPath := filepath.Join(configDir, fmt.Sprintf("%s_config.json", commUnitID))
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Config file not found"})
			return
		}

		configData, err := os.ReadFile(configPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read config file"})
			return
		}

		var config interface{}
		if err := json.Unmarshal(configData, &config); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid JSON in config file"})
			return
		}

		c.JSON(http.StatusOK, config)
	})

	router.POST("/postData", func(c *gin.Context) {
		var data map[string]interface{}
		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		fmt.Printf("Received data: %v\n", data)

		f, err := os.OpenFile("received_data.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write data"})
			return
		}
		defer f.Close()

		dataBytes, _ := json.Marshal(data)
		f.WriteString(string(dataBytes) + "\n")

		c.JSON(http.StatusOK, gin.H{"status": "Data received"})
	})

	router.Run(":5000")
}
