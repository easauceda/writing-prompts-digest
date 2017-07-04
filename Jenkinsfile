pipeline {
  agent {
    docker {
      image 'golang'
    }
    
  }
  stages {
    stage('Get') {
      steps {
        sh 'go get github.com/easauceda/writing-prompts-digest'
      }
    }
    stage('Build') {
      steps {
        sh 'go build *.go'
      }
    }
  }
}