pipeline {
    agent any

    environment {
        APP_NAME = "daily-coffee-api"
        DOCKER_TAG = "local"
        DOCKER_COMPOSE_FILE = 'docker-compose.yml'
    }

    stages {
        stage('Checkout') {
            steps {
                checkout([$class: 'GitSCM', 
                    branches: [[name: '*/master']],
                    userRemoteConfigs: [[url: 'https://github.com/ojihalaw/ojihalaw-daily-coffee-api.git']],
                    extensions: [[$class: 'CleanBeforeCheckout']]
                ])
            }
        }


        stage('Build Go App') {
            steps {
                withCredentials([file(credentialsId: 'daily-coffee-env', variable: 'ENV_FILE')]) {
                    sh 'cp $ENV_FILE .env'
                    sh 'docker-compose build'
                }
            }
        }

        stage('Deploy') {
            steps {
                sh '''
                docker-compose down
                docker-compose up -d
                '''
            }
        }
    }
}
