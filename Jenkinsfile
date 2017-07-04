pipeline {
  agent {
    docker {
      image 'golang'
    }
    
  }
  stages {
    stage('Get') {
      steps {
        sh 'go get -v github.com/easauceda/writing-prompts-digest'
      }
    }
    stage('Build') {
      steps {
        sh 'go build *.go'
      }
    }
    stage('Test') {
      steps {
        echo 'testing'
      }
    }
    stage('Deploy') {
      steps {
        echo 'Deploying!'
      }
    }
  }
}