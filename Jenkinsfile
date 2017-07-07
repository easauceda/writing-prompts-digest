pipeline {
  agent any
  environment {
    GIT_SHA = sh "git rev-parse HEAD" 
  }
  stages {
    stage('Build') {
      steps {
        sh 'printenv'
        sh "docker build -t quay.io/easauceda/writing-prompts-digest:`git rev-parse HEAD`."
      }
    }
    stage('Deploy') {
      steps {
        withCredentials([usernamePassword(credentialsId: 'quay_credentials', passwordVariable: 'quay_pw', usernameVariable: 'quay_username')]) {
          sh "docker login -u=${quay_username} -p=${quay_pw} quay.io"
          sh "docker push quay.io/easauceda/writing-prompts-digest:`git rev-parse HEAD`"
        }
        echo 'Deploying'
      }
    }
  }
}
