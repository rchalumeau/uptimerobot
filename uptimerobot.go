package main
 
import (
    "fmt"
    "strings"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "time"
    "strconv"
)
// Structures from uptime robot
// URL being monitored
type Server struct {
   Name    string  `json:"friendly_name"`
   Url     string
   Custom_Uptime_Ranges   string  "custom_uptime_ranges"
}
// Overall structure
type Monitoring struct {
   Monitors []Server
}
// Returns the epoch difference to be sent as para,eter to uptime
// to get the data for a month
func getDates(increment int) (string, string) {
   target := time.Now().AddDate(0, -increment, 0)
   end_target := target.AddDate(0, 1, 0)
   first_day := time.Date(target.Year(), target.Month(), 1, 0, 0, 0, 0, time.UTC)
   last_day := time.Date(end_target.Year(), end_target.Month(), 1, 0, 0, 0, 0, time.UTC)
   epoch := strconv.FormatInt(first_day.Unix(), 10)+ "_" + strconv.FormatInt(last_day.Unix(), 10)
   month_string := fmt.Sprintf("%s %d", target.Month(), target.Year())
   return epoch, month_string
} 
var url = "https://api.uptimerobot.com/v2/getMonitors"
func getData(epoch string, api_key string, monitors string) Monitoring {
    var m Monitoring
    // Preparation of the HTTP request to uptimerobot API
    parameters := fmt.Sprintf("api_key=%s&format=json&logs=1&monitors=%s&custom_uptime_ranges=%s", api_key, monitors, epoch)
    //fmt.Println(parameters, "\n")
    payload := strings.NewReader(parameters)
    req, _ := http.NewRequest("POST", url, payload)
    req.Header.Add("content-type", "application/x-www-form-urlencoded")
    req.Header.Add("cache-control", "no-cache")
    res, _ := http.DefaultClient.Do(req)
    defer res.Body.Close()
    body, _ := ioutil.ReadAll(res.Body)
    //fmt.Println("response:", string(body))
    // Json parsing in structures
    err := json.Unmarshal(body, &m)
    if err != nil {
        fmt.Println("error:", err)
    }
    return m
}
func main() {
    // Variables to access uptimerobot
    api_key := "u473525-4dd79869f3484bf704cc6f89" 
    monitors := "779311773-779409798-780664703"
    number_of_months := 3
    for i := 1; i <= number_of_months; i++ {
        epoch, month := getDates(i)
        m := getData(epoch, api_key, monitors)
        fmt.Println(month)
        for _, server := range m.Monitors {
            fmt.Printf("\t%s (%s) : %v %%\n", server.Name, server.Url, server.Custom_Uptime_Ranges )
        }
    }
}