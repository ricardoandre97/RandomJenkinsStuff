package main

import (
    "log"
    "text/template"
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "bytes"
    "fmt"
)

type AutoJobs struct {
    Jobs    []Job    `yaml:"jobs"`
    Folders []Folder `yaml:"folders"`
}

type Job struct {
    JobName    string     `yaml:"job_name"`
    JobDesc    string     `yaml:"job_description"`
    Params     []JobParam `yaml:"parameters"`
    GitURL     string     `yaml:"git_url"`
    CredsID    string     `yaml:"creds_id"`
    Branch     string     `yaml:"branch"`
    ScriptPath string     `yaml:"script_path"`
    JobFolder  string     `yaml:"folder"`
}

type JobParam struct {
    Name   string  `yaml:"name"`
    Value  string  `yaml:"value"`
    Desc   string  `yaml:"desc"`
    Type   string  `yaml:"type"`
}

type Folder struct {
    FolderName string `yaml:"folder_name"`
    FolderDesc string `yaml:"folder_desc"`
}

func getDSLTemplate(file string) (string, error) {

    // Read jobs file
    data, dErr := ioutil.ReadFile(file)
    if dErr != nil {
        return "", dErr
    }

    // Unmasrshal to struct
    automatedJobs := AutoJobs{}
    err := yaml.Unmarshal(data, &automatedJobs)
    if err != nil {
        return "", err
    }

    // Parse template
    t := template.Must(template.ParseFiles("templates/dsl.tpl"))

    var byteTpl bytes.Buffer
    if err := t.Execute(&byteTpl, automatedJobs); err != nil {
        return "", err
    }

    return byteTpl.String(), nil

}


func getDSLJobs(dir string) (string, error) {

    var r string

    files, err := ioutil.ReadDir(dir)
    if err != nil {
        return r, err
    }

    
    for _, f := range files {
        str, err := getDSLTemplate(dir+"/"+f.Name())
        if err != nil {
            return r, err
        }
        r += str+"\n"
    }

    return r, nil

}

func main() {
    dsl, err := getDSLJobs("./jobs")
    if err != nil {
        log.Fatalf("Error rendering jobs: %v", err)
    }
    fmt.Println(dsl)
}