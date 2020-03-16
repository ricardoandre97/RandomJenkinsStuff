package main

import (
    "log"
    "os"
    "html/template"
    "gopkg.in/yaml.v2"
    "io/ioutil"
)

type Job struct {
    JobName    string  `yaml:"job_name"`
    JobDesc    string  `yaml:"job_description"`
    Params     []Param `yaml:"parameters"`
    GitURL     string  `yaml:"git_url"`
    CredsID    string  `yaml:"creds_id"`
    Branch     string  `yaml:"branch"`
    ScriptPath string  `yaml:"script_path"`
}

type Param struct {
    Name   string  `yaml:"name"`
    Value  string  `yaml:"value"`
    Desc   string  `yaml:"desc"`
    Type   string  `yaml:"type"`
}

func main() {

    // Read jobs file
    data, dErr := ioutil.ReadFile("jobs.yaml")
    if dErr != nil {
        log.Fatalf("File reading error", dErr)
    }

    // Unmasrshal to struct
    jobs := []Job{}
    err := yaml.Unmarshal(data, &jobs)
    if err != nil {
        log.Fatalf("error: %v", err)
    }

    // Parse template
    tmpl := template.Must(template.ParseFiles("templates/dsl.tpl"))

    // Create the file
    f, err := os.Create("job.dsl")
    if err != nil {
        log.Fatalf("error: %v", err)
    }

    // Execute the template to the file.
    err = tmpl.Execute(f, jobs)
    if err != nil {
        log.Fatalf("error: %v", err)
    }

    // Close the file when done.
    f.Close()
}

