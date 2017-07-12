pipeline {
  agent any
  stages {
    stage('Test') {
      agent {
        docker {
          image 'golang'
        }
      }
      steps {
        sh "go get -v -t -d ./..."
        sh "go test"
      }
    }
    stage('Deploy') {
      when {
        branch 'master'
      }
      steps {
        sh "docker build -t quay.io/easauceda/writing-prompts-digest:`git rev-parse HEAD` ."
        withCredentials([usernamePassword(credentialsId: 'quay_credentials', passwordVariable: 'quay_pw', usernameVariable: 'quay_username')]) {
          sh "docker login -u=${quay_username} -p=${quay_pw} quay.io"
          sh "docker push quay.io/easauceda/writing-prompts-digest:`git rev-parse HEAD`"
        }
        echo 'Deploying'
      }
    }
  }
}
