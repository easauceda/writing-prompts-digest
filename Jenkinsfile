pipeline {
  agent any
  stages {
    stage('Build') {
      steps {
        sh 'env'
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
