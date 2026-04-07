package services

import (
    "bytes"
    "encoding/json"
    "log"
    "math/rand"
    "net/http"
    "time"

    "monitoring-backend/models"
)

type Agent struct {
    NodeName string
    APIURL   string
    Token    string
}

func NewAgent(nodeName, apiURL string) *Agent {
    return &Agent{
        NodeName: nodeName,
        APIURL:   apiURL,
    }
}

func (a *Agent) Start(intervalSeconds int) {
    if err := a.login(); err != nil {
        log.Fatalf("Agent failed to login: %v", err)
    }

    ticker := time.NewTicker(time.Duration(intervalSeconds) * time.Second)
    for range ticker.C {
        a.reportMetrics()
    }
}

func (a *Agent) login() error {
    req := models.AuthRequest{
        Username: "admin",
        Password: "admin123",
    }
    
    body, _ := json.Marshal(req)
    resp, err := http.Post(a.APIURL+"/auth/login", "application/json", bytes.NewBuffer(body))
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    var authRes models.AuthResponse
    if err := json.NewDecoder(resp.Body).Decode(&authRes); err != nil {
        return err
    }

    a.Token = authRes.Token
    log.Println("Agent logged in successfully")
    return nil
}

func (a *Agent) reportMetrics() {
    metric := models.Metric{
        Time:        time.Now(),
        NodeName:    a.NodeName,
        CPUUsage:    rand.Float64() * 100,
        MemoryUsage: rand.Float64() * 100,
        DiskUsage:   rand.Float64() * 100,
    }

    body, _ := json.Marshal(metric)
    req, _ := http.NewRequest("POST", a.APIURL+"/metrics", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+a.Token)

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        log.Printf("Failed to report metrics: %v", err)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusCreated {
        log.Printf("Reported metrics for %s", a.NodeName)
    } else {
        log.Printf("Failed to report metrics, status: %d", resp.StatusCode)
    }
}
