pipeline {
    agent any
    stages {
        stage('Create DSL') {
            agent {
                docker { image 'jobcreator:v1'}
            }
            steps {
                sh 'jobCreator' // This is a binary available in PATH
            }
        }
        stage('Create Jobs') {
            agent {
                docker { image 'jobcreator:v1'}
            }
            steps {
                jobDsl removedJobAction: 'DELETE', targets: 'job.dsl'
            }
        }
    }
}