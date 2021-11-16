#!/usr/bin/env groovy
// vim: ft=Jenkinsfile

node('go1.17') {
	container('run'){
		def newTag = ''
		def tag = ''
		def gitTag = ''

		try {
			stage('Checkout'){
				checkout scm
				gitTag = sh(script: 'git tag -l --contains HEAD', returnStdout: true).trim()
			}

			stage('Fetch dependencies'){
				// using ID because: https://issues.jenkins-ci.org/browse/JENKINS-32101
				sh('go mod download')
			}
			stage('Run test'){
				sh('make test')
			}

			if(gitTag != ''){
				tag = gitTag
			}else if (env.BRANCH_NAME == 'main'){
				echo "Skipping build since we did not find an existing tag"
				return
			}

			if( tag != ''){
				stage('Build the application'){
					echo "Building with docker tag ${tag}"
					docker.withRegistry("https://quay.io", 'docker-registry') {
						sh("VERSION=${tag} make push")
					}
				}
			}

			currentBuild.result = 'SUCCESS'
		} catch(err) {
			currentBuild.result = 'FAILED'
			throw err
		}
	}
}
