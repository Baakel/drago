package main

import (
    "encoding/json"
    "github.com/joho/godotenv"
    "image"
    _ "image/jpeg"
    "io/ioutil"
    "log"
    "math"
    "math/rand"
    "net/http"
    "os"
    "strconv"
    "time"
)

func getRandomPage() int {
    var m map[string]interface{}
    wallHavenUrl := getEnvVar("WALLHAVEN_URL")
    resp, err := http.Get(wallHavenUrl)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    jsonObjects, err := ioutil.ReadAll(resp.Body)

    if err := json.Unmarshal(jsonObjects, &m); err != nil {
        log.Fatal(err)
    }
    entries, _ := m["data"].([]interface{})
    for _, entry := range entries {
        entryMap := entry.(map[string]interface{})
        if entryMap["label"] == "reference" {
            images := entryMap["count"].(float64)
            pages := int(RoundToNextInt(images / 24))
            s1 := rand.NewSource(time.Now().UnixNano())
            r1 := rand.New(s1)
            return r1.Intn(pages) + 1
        }
    }
    return 1
}

func RoundToNextInt(x float64) float64 {
    t := math.Trunc(x)
    if math.Abs(x-t) > 0 {
        return t + 1
    }
    return t
}

func getRandomPath(page int) string {
    wallHavenUrl := getEnvVar("WALLHAVEN_URL")
    apiKey := getEnvVar("API_KEY")
    var m map[string]interface{}
    resp, err := http.Get(wallHavenUrl + apiKey + strconv.Itoa(page))
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    jsonObjects, err := ioutil.ReadAll(resp.Body)

    if err := json.Unmarshal(jsonObjects, &m); err != nil {
        log.Fatal(err)
    }

    entry, _ := m["data"].([]interface{})
    entrySize := len(entry)
    s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)

    object := entry[r1.Intn(entrySize)]
    objectMap := object.(map[string]interface{})

    return objectMap["path"].(string)
}

func getWall() image.Image {
    resp, err := http.Get(getRandomPath(getRandomPage()))
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    img, _, err := image.Decode(resp.Body)
    if err != nil {
        log.Fatal(err)
    }

    return img
}

func getEnvVar(key string) string {
    err := godotenv.Load()

    if err != nil {
        log.Fatalf("Error loading env file")
    }

    return os.Getenv(key)
}
