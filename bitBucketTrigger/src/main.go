package main

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "errors"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "strings"
)

type jenkins struct {
    Host     string
    User     string
    Password string
}

type gitBody struct {
    Repository struct {
        Name string `json:"name"`
    }
}

func requestIsValid(body []byte, secret string, gitHeader string) bool {
    if gitHeader == "" {
        log.Print("Can't validate request since X-Hub-Signature header is empty")
        return false
    }

    // New hmac with 256
    hash := hmac.New(sha256.New, []byte(secret))
    _, err := hash.Write(body)

    if err != nil {
        log.Printf("Couldn't write body to hash: %v", err)
        return false
    }

    // Create the expected signature
    signature := hex.EncodeToString(hash.Sum(nil))

    // Remove the "=" from the Git header and compare
    return signature == strings.Split(gitHeader, "=")[1]
}

func triggerJenkins(j jenkins, job string) error {
    // Get the Crumb from Jenkins
    crumbURL := fmt.Sprintf("%v/crumbIssuer/api/xml?xpath=concat(//crumbRequestField,\":\",//crumb)", j.Host)
    jenkinsURL := fmt.Sprintf("%v/job/%v/build?delay=0", j.Host, job) // This should be dynamic

    // New client with request
    client := &http.Client{}
    req, err := http.NewRequest("GET", crumbURL, nil)

    // Make request to get crumb
    req.SetBasicAuth(j.User, j.Password)
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // Load crumb value
    crumb, err := ioutil.ReadAll(resp.Body)
    header := strings.Split(string(crumb), ":")

    if err != nil {
        return err
    }
    if resp.StatusCode != 200 {
        return errors.New(fmt.Sprintf("Expected 200 getting crumb, but got %v", resp.StatusCode))
    }

    // Start the POST request to trigger Jenkins
    // This needs to be improved to reuse code
    postReq, postErr := http.NewRequest("POST", jenkinsURL, nil)
    postReq.SetBasicAuth(j.User, j.Password)
    postReq.Header.Set(header[0], header[1])
    postResp, postErr := client.Do(postReq)
    if postErr != nil {
        return postErr
    }
    defer postResp.Body.Close()

    jResp, jErr := ioutil.ReadAll(postResp.Body)
    log.Printf("Tiggering job gave %v", postResp.StatusCode)
    if jErr != nil {
        return jErr
    }

    if !(postResp.StatusCode >= 200 && postResp.StatusCode <= 299) {
        log.Print(string(jResp))
        return errors.New(fmt.Sprintf("Expected 200-299 triggering job but got %v", postResp.StatusCode))
    }
    return nil

}

func handler(w http.ResponseWriter, r *http.Request) {

    // Methos should only be post
    if r.Method != "POST" {
        w.WriteHeader(401)
        return
    }

    // Secret env must be there to validate the request
    secret := os.Getenv("SECRET")

    if secret == "" {
        log.Printf("SECRET not present")
        w.WriteHeader(400)
        return
    }

    // Get the body payload to check the request
    body, bodyErr := ioutil.ReadAll(r.Body)

    if bodyErr != nil {
        log.Printf("Couldn't read incoming body request: %v", bodyErr)
        w.WriteHeader(500)
        return 
    }
    if !requestIsValid(body, secret, r.Header.Get("X-Hub-Signature")) {
        log.Printf("Request rejected. Signatures don't match")
        w.WriteHeader(400)
        w.Write([]byte("Request rejected. Signatures don't match"))
        return
    }

    // If request is valid, check that we have all we need to trigger jobs
    req := jenkins{
        Host:     os.Getenv("JENKINS_HOST"),
        User:     os.Getenv("API_USER"),
        Password: os.Getenv("API_PASSWORD"),
    }

    if req.Host == "" || req.User == "" || req.Password == "" {
        log.Printf("JENKINS_HOST, API_USER, API_PASSWORD vars are mandatory")
        w.WriteHeader(400)
        return
    }

    // Get the name of the repo
    goBody := gitBody{}
    jsonErr := json.Unmarshal(body, &goBody)

    if jsonErr != nil {
        log.Printf("Coulnd't get repo name from bitbucket: %v", jsonErr)
        w.WriteHeader(500)
        return  
    }

    // Trigger Pipeline
    err := triggerJenkins(req, goBody.Repository.Name)
    if err != nil {
        msg := fmt.Sprintf("Couldn't trigger Jenkins job: %v", err)
        log.Printf(msg)
        w.WriteHeader(400)
        w.Write([]byte(msg))
        return
    }
    log.Println("Ok")
    w.WriteHeader(200)
}

func main() {
    http.HandleFunc("/gw", handler)
    http.ListenAndServe(":9090", nil)
}