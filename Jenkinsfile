pipeline {
  agent none
  stages {
    stage('Build') {
      steps {
        sh 'curl -o ca-certificates.crt https://raw.githubusercontent.com/bagder/ca-bundle/master/ca-bundle.crt'
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