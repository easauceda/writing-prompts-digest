pipeline {
  agent any
  stages {
    stage('Build') {
      steps {
        sh "docker build -t writing-prompts-digest:`git rev-parse --short HEAD` ."
      }
    }
    stage('Deploy') {
      steps {
        echo 'Deploying'
      }
    }
  }
}
