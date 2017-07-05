pipeline {
  agent none
  stages {
    stage('Build') {
      steps {
        sh 'docker build -t wpd .'
      }
    }
    stage('Deploy') {
      steps {
        echo 'Deploying!'
      }
    }
  }
}