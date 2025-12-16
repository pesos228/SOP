pipeline {
    agent any
    environment {
        COMPOSE_PROJECT_NAME = "sop"
    }
    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }
        stage('Test (CI)') {
            steps {
                script {
                    echo "Запуск тестов..."
                    
                    sh """
                        docker run --rm \
                        -v "\$PWD":/workspace \
                        -v go-mod-cache:/go/pkg/mod \
                        -v go-build-cache:/root/.cache/go-build \
                        -w /workspace/hosting-service \
                        golang:1.24.4-alpine \
                        go test -v ./...
                    """
                }
            }
        }
        stage('Deploy (CD)') {
            steps {
                script {
                    echo "Тесты прошли успешно. Деплоим..."
                    
                    sh '''
                        if command -v docker-compose >/dev/null 2>&1; then
                            docker-compose up -d --build --no-deps hosting-service provisioning-service migrator
                        else
                            docker compose up -d --build --no-deps hosting-service provisioning-service migrator
                        fi
                    '''
                }
            }
        }
        stage('Health Check') {
            steps {
                script {
                    sleep 10
                    sh 'curl -f http://hosting-api:8080/health || echo "Health check failed (hosting-api)"'
                    sh 'curl -f http://hosting-provisioner:7070/health || echo "Health check failed (hosting-provisioner)"'
                }
            }
        }
    }
}