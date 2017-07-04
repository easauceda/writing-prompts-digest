pipeline {
  agent {
    docker {
      image 'golang'
    }
    
  }
  stages {
    stage('Hello World!') {
      steps {
        sh 'go get github.com/easauceda/writing-prompts-digest'
      }
    }
  }
}