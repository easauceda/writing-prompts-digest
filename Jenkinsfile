pipeline {
  agent any
  stages {
    stage('Build') {
      steps {
        sh "echo $(git rev-parse --short HEAD)"
      }
    }
    stage('Deploy') {
      steps {
        echo 'Deploying'
      }
    }
  }
}
