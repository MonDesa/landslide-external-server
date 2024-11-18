package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "time"

    MQTT "github.com/eclipse/paho.mqtt.golang"
    "github.com/gin-gonic/gin"
)

type DeviceInfo struct {
    CommUnitID string `json:"CommUnitID"`
}

var devices map[string]DeviceInfo
var devicesFile = "devices.json"
var mqttClient MQTT.Client

func main() {
    router := gin.Default()

    configDir := "./configs"
    if _, err := os.Stat(configDir); os.IsNotExist(err) {
        os.Mkdir(configDir, os.ModePerm)
    }

    devices = make(map[string]DeviceInfo)
    loadDevices()

    initMQTTClient()

    router.POST("/configs/:commUnitID", updateConfigHandler)

    router.Run(":5000")
}

func initMQTTClient() {
    opts := MQTT.NewClientOptions()
    opts.AddBroker("tcp://your_mqtt_broker_ip:1883")
    opts.SetClientID("backend_server")
    opts.SetUsername("your_mqtt_username")
    opts.SetPassword("your_mqtt_password")

    mqttClient = MQTT.NewClient(opts)
    if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
        panic(token.Error())
    }

    fmt.Println("Connected to MQTT broker")
}

func loadDevices() {
    if _, err := os.Stat(devicesFile); err == nil {
        data, err := ioutil.ReadFile(devicesFile)
        if err == nil {
            json.Unmarshal(data, &devices)
        }
    }
}

func saveDevices() error {
    devicesBytes, err := json.MarshalIndent(devices, "", "  ")
    if err != nil {
        return err
    }

    tempDevicesFile := devicesFile + ".tmp"
    if err := ioutil.WriteFile(tempDevicesFile, devicesBytes, 0644); err != nil {
        return err
    }

    return os.Rename(tempDevicesFile, devicesFile)
}

func updateConfigHandler(c *gin.Context) {
    commUnitID := c.Param("commUnitID")

    var configData map[string]interface{}
    if err := c.BindJSON(&configData); err != nil {
        c.JSON(400, gin.H{"error": "Invalid JSON"})
        return
    }

    configPath := filepath.Join("./configs", fmt.Sprintf("%s_config.json", commUnitID))
    configBytes, err := json.MarshalIndent(configData, "", "  ")
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to marshal config data"})
        return
    }

    tempConfigPath := configPath + ".tmp"
    if err := ioutil.WriteFile(tempConfigPath, configBytes, 0644); err != nil {
        c.JSON(500, gin.H{"error": "Failed to write temp config file"})
        return
    }

    if err := os.Rename(tempConfigPath, configPath); err != nil {
        c.JSON(500, gin.H{"error": "Failed to save config file"})
        return
    }

    topic := fmt.Sprintf("comm_unit/%s/config", commUnitID)
    token := mqttClient.Publish(topic, 1, false, configBytes)
    token.Wait()
    if token.Error() != nil {
        c.JSON(500, gin.H{"error": "Failed to publish config via MQTT"})
        return
    }

    fmt.Printf("Published config to topic %s\n", topic)

    statusTopic := fmt.Sprintf("comm_unit/%s/config/status", commUnitID)
    statusChan := make(chan string)
    mqttClient.Subscribe(statusTopic, 1, func(client MQTT.Client, msg MQTT.Message) {
        statusChan <- string(msg.Payload())
    })

    select {
    case statusMsg := <-statusChan:
        fmt.Printf("Received status from device %s: %s\n", commUnitID, statusMsg)
        c.JSON(200, gin.H{"status": "Config updated and applied by device", "deviceStatus": statusMsg})
    case <-time.After(10 * time.Second):
        c.JSON(500, gin.H{"error": "Timeout waiting for device status"})
    }

    mqttClient.Unsubscribe(statusTopic)
}