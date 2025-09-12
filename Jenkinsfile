pipeline {
    agent any

    environment {
        APP_NAME = "daily-coffee-api"
        DOCKER_TAG = "local"
        DOCKER_COMPOSE_FILE = 'docker-compose.yml'
    }

    stages {
        stage('Cleanup Workspace') {
            steps {
                // Clean workspace before starting
                deleteDir()
            }
        }

        stage('Checkout') {
            steps {
                script {
                    try {
                        // Method 1: Using checkout scm (recommended for pipeline jobs)
                        checkout scm
                    } catch (Exception e) {
                        echo "SCM checkout failed, trying explicit git checkout..."
                        // Method 2: Explicit git checkout
                        checkout([$class: 'GitSCM', 
                            branches: [[name: '*/master']],
                            userRemoteConfigs: [[url: 'https://github.com/ojihalaw/ojihalaw-daily-coffee-api.git']],
                            extensions: [
                                [$class: 'CleanBeforeCheckout'],
                                [$class: 'CloneOption', depth: 1, noTags: false, reference: '', shallow: true]
                            ]
                        ])
                    }
                }
                
                // Verify checkout worked
                sh 'pwd && ls -la'
                sh 'git status || echo "Not a git repository - that\'s okay for some operations"'
            }
        }

        stage('Verify Environment') {
            steps {
                sh '''
                echo "Current directory: $(pwd)"
                echo "Directory contents:"
                ls -la
                echo "Docker version:"
                docker --version
                echo "Docker Compose version:"
                docker compose --version || docker compose version
                '''
            }
        }

        stage('Build Go App') {
            steps {
                withCredentials([file(credentialsId: 'daily-coffee-env', variable: 'ENV_FILE')]) {
                    sh '''
                    echo "Copying environment file..."
                    cp $ENV_FILE .env
                    
                    echo "Checking if docker-compose.yml exists..."
                    if [ ! -f docker-compose.yml ]; then
                        echo "docker-compose.yml not found!"
                        exit 1
                    fi
                    
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
                docker compose down --remove-orphans
                
                echo "Starting new containers..."
                docker compose up -d
                
                echo "Checking container status..."
                docker compose ps
                '''
            }
        }

        stage('Health Check') {
            steps {
                script {
                    // Wait for the application to start
                    sleep(time: 30, unit: 'SECONDS')
                    
                    // Add your health check here
                    sh '''
                    echo "Performing health check..."
                    docker-compose ps
                    
                    # Uncomment and modify the following line for HTTP health check
                    # curl -f http://localhost:8080/health || exit 1
                    '''
                }
            }
        }
    }

    post {
        always {
            echo 'Pipeline completed!'
            // Clean up if needed
            sh 'docker system prune -f --volumes || true'
        }
        success {
            echo 'Pipeline succeeded!'
        }
        failure {
            echo 'Pipeline failed!'
            sh 'docker-compose logs || true'
        }
    }
}