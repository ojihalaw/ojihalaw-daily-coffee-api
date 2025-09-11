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
                git branch: 'main', url: 'https://github.com/ojihalaw/ojihalaw-daily-coffee-api'
            }
        }
        
        stage('Build Go App') {
            steps {
                sh 'docker-compose build'
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
