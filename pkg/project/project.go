package project

var (
	description        = "aws-collector is a prometheus exporter for AWS Control Planes"
	gitSHA             = "n/a"
	name        string = "aws-collector"
	source      string = "https://github.com/giantswarm/aws-collector"
	version            = "1.1.0"
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
