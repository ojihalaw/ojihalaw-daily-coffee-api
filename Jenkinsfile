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
                // Clean workspace and checkout code
                deleteDir()
                sh '''
                git clone https://github.com/ojihalaw/ojihalaw-daily-coffee-api.git .
                git checkout master
                echo "Repository cloned successfully"
                ls -la
                '''
            }
        }

        stage('Verify Environment') {
            steps {
                sh '''
                echo "Current directory: $(pwd)"
                echo "Directory contents:"
                ls -la
                echo "Checking for docker-compose.yml:"
                test -f docker-compose.yml && echo "docker-compose.yml found" || echo "docker-compose.yml NOT found"
                '''
            }
        }

        stage('Build Go App') {
            steps {
                withCredentials([file(credentialsId: 'daily-coffee-env', variable: 'ENV_FILE')]) {
                    sh '''
                    echo "Copying environment file..."
                    cp "$ENV_FILE" .env
                    
                    echo "Checking Docker installation..."
                    which docker || echo "Docker not found in PATH"
                    docker --version || echo "Docker command failed"
                    
                    echo "Checking Docker Compose installation..."
                    which docker-compose || echo "docker-compose not found"
                    docker-compose --version || echo "docker-compose command failed"
                    
                    echo "Checking if Docker daemon is accessible..."
                    docker info || echo "Cannot connect to Docker daemon"
                    
                    echo "Building with Docker Compose..."
                    docker-compose build
                    '''
                }
            }
        }

        stage('Deploy') {
            steps {
                sh '''
                echo "Stopping existing containers..."
                docker-compose down --remove-orphans || true
                
                echo "Starting new containers..."
                docker-compose up -d
                
                echo "Checking container status..."
                docker-compose ps
                '''
            }
        }
    }

    post {
        always {
            echo 'Cleaning up...'
            sh 'docker system prune -f || true'
        }
        failure {
            echo 'Build failed! Checking logs...'
            sh 'docker-compose logs || true'
        }
    }
}