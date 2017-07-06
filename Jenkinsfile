pipeline {
  agent any
  stages {
    stage('Build') {
      steps {
        sh "echo ${git rev-parse --short HEAD}"
        //sh 'docker build -t wpd:$(GIT_COMMIT) .'
      }
    }
    stage('Deploy') {
      steps {
        echo 'Deploying'
      }
    }
  }
}
