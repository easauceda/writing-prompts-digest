pipeline {
  agent any
  stages {
    stage('Build') {
      steps {
        sh 'docker build -t wpd:$(GIT_COMMIT) .'
      }
    }
    stage('Deploy') {
      steps {
        echo 'Deploying'
      }
    }
  }
}
