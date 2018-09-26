pipeline {
    agent any

    stages {
        stage('Checkout')
        {
            steps{
                checkout scm
            }
        }

        stage('Build') {
            steps {
                echo 'Building..'
                sh 'make docker'
                junit '**/build/*.xml'
            }
        }
   }

    post {
        always{
            cleanWs()
        }
    }
}
