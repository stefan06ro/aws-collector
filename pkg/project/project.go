package project

var (
	description        = "The aws-collector manages Kubernetes clusters running on AWS."
	gitSHA             = "n/a"
	name        string = "aws-collector"
	source      string = "https://github.com/giantswarm/aws-collector"
	version            = "0.0.1-dev"
)

func Description() string {
	return description
}

func GitSHA() string {
	return gitSHA
}

func Name() string {
	return name
}

func Source() string {
	return source
}

func Version() string {
	return version
}
