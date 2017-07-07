pipeline {
  agent any
  stages {
    stage('Build') {
      steps {
        sh "export GIT_SHA=$(git rev-parse HEAD)"
        sh "docker build -t quay.io/easauceda/writing-prompts-digest:$GIT_SHA."
      }
    }
    stage('Deploy') {
      steps {
        withCredentials([usernamePassword(credentialsId: 'quay_credentials', passwordVariable: 'quay_pw', usernameVariable: 'quay_username')]) {
          sh "docker login -u=${quay_username} -p=${quay_pw} quay.io"
          sh "docker push quay.io/easauceda/writing-prompts-digest:$GIT_SHA"
        }
        echo 'Deploying'
      }
    }
  }
}
