pipeline {
  agent any
  environment {
    GIT_SHA = `git rev-parse --short HEAD`
  }
  stages {
    stage('Build') {
      steps {
        sh "docker build -t quay.io/easauceda/writing-prompts-digest:${env.GIT_SHA}."
      }
    }
    stage('Deploy') {
      steps {
        withCredentials([usernamePassword(credentialsId: 'quay_credentials', passwordVariable: 'quay_pw', usernameVariable: 'quay_username')]) {
          "docker login -u=${quay_username} -p=${quay_pw} quay.io"
          "docker push quay.io/easauceda/writing-prompts-digest:${env.GIT_SHA}"
        }
        echo 'Deploying'
      }
    }
  }
}
