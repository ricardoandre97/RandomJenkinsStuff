## Bitbucket & Jenkins
# Install the bitBucket server notifier plugin and configure it on Jenkins
# Install strict issuer plugin and configure it
# Create a MultiPipeline Job + ssh keys in BitBucket -> Job name should be repo name
# Point the bitBucket webhook to http://{{jenkins-host:9090}}/gw (make sure this is running)
# Put this Jenkinsfile in your repo
pipeline {
    agent any
    stages {
        stage('Build') {
            steps {
                sh 'echo build'
            }
        }
        stage('Test') {
            steps {
                sh 'echo test'
            }
        }
        stage('sonar analytics') {
            when {
                expression { BRANCH_NAME ==~ /(master|development)/ }  
            }
            steps {
                sh 'echo Sonar'
            }
        }
        stage('DEB stage') {
            when {
                branch 'master'
            }
            steps {
                 sh 'echo Building DEB'
            }
        }
    }
    post {
        always {
            script {
                currentBuild.result = currentBuild.result ?: 'SUCCESS'
                notifyBitbucket()
            }
        }
    }
}