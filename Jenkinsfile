pipeline {
  agent {
    docker {
      image 'busybox'
    }
    
  }
  stages {
    stage('Hello World!') {
      steps {
        sh 'echo hi'
      }
    }
  }
}