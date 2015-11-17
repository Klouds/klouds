package models

type MarathonApplication struct {
	App struct {
		AcceptedResourceRoles interface{} `json:"acceptedResourceRoles"`
		Args                  []string    `json:"args"`
		BackoffFactor         float64     `json:"backoffFactor"`
		BackoffSeconds        int         `json:"backoffSeconds"`
		Cmd                   interface{} `json:"cmd"`
		Constraints           [][]string  `json:"constraints"`
		Container             struct {
			Docker struct {
				ForcePullImage bool          `json:"forcePullImage"`
				Image          string        `json:"image"`
				Network        string        `json:"network"`
				Parameters     []interface{} `json:"parameters"`
				PortMappings   []struct {
					ContainerPort int    `json:"containerPort"`
					HostPort      int    `json:"hostPort"`
					Protocol      string `json:"protocol"`
					ServicePort   int    `json:"servicePort"`
				} `json:"portMappings"`
				Privileged bool `json:"privileged"`
			} `json:"docker"`
			Type    string        `json:"type"`
			Volumes []interface{} `json:"volumes"`
		} `json:"container"`
		Cpus         float64       `json:"cpus"`
		Dependencies []interface{} `json:"dependencies"`
		Deployments  []interface{} `json:"deployments"`
		Disk         int           `json:"disk"`
		Env          struct{}      `json:"env"`
		Executor     string        `json:"executor"`
		HealthChecks []struct {
			GracePeriodSeconds     int    `json:"gracePeriodSeconds"`
			IgnoreHTTP1xx          bool   `json:"ignoreHttp1xx"`
			IntervalSeconds        int    `json:"intervalSeconds"`
			MaxConsecutiveFailures int    `json:"maxConsecutiveFailures"`
			Path                   string `json:"path"`
			PortIndex              int    `json:"portIndex"`
			Protocol               string `json:"protocol"`
			TimeoutSeconds         int    `json:"timeoutSeconds"`
		} `json:"healthChecks"`
		ID        string `json:"id"`
		Instances int    `json:"instances"`
		Labels    struct {
			HAPROXYHTTP      string `json:"HAPROXY_HTTP"`
			HTTPPORTIDX0NAME string `json:"HTTP_PORT_IDX_0_NAME"`
		} `json:"labels"`
		MaxLaunchDelaySeconds int           `json:"maxLaunchDelaySeconds"`
		Mem                   int           `json:"mem"`
		Ports                 []int         `json:"ports"`
		RequirePorts          bool          `json:"requirePorts"`
		StoreUrls             []interface{} `json:"storeUrls"`
		Tasks                 []struct {
			AppID              string `json:"appId"`
			HealthCheckResults []struct {
				Alive               bool        `json:"alive"`
				ConsecutiveFailures int         `json:"consecutiveFailures"`
				FirstSuccess        string      `json:"firstSuccess"`
				LastFailure         interface{} `json:"lastFailure"`
				LastSuccess         string      `json:"lastSuccess"`
				TaskID              string      `json:"taskId"`
			} `json:"healthCheckResults"`
			Host      string `json:"host"`
			ID        string `json:"id"`
			Ports     []int  `json:"ports"`
			SlaveID   string `json:"slaveId"`
			StagedAt  string `json:"stagedAt"`
			StartedAt string `json:"startedAt"`
			Version   string `json:"version"`
		} `json:"tasks"`
		TasksHealthy    int `json:"tasksHealthy"`
		TasksRunning    int `json:"tasksRunning"`
		TasksStaged     int `json:"tasksStaged"`
		TasksUnhealthy  int `json:"tasksUnhealthy"`
		UpgradeStrategy struct {
			MaximumOverCapacity   int `json:"maximumOverCapacity"`
			MinimumHealthCapacity int `json:"minimumHealthCapacity"`
		} `json:"upgradeStrategy"`
		Uris        []interface{} `json:"uris"`
		User        interface{}   `json:"user"`
		Version     string        `json:"version"`
		VersionInfo struct {
			LastConfigChangeAt string `json:"lastConfigChangeAt"`
			LastScalingAt      string `json:"lastScalingAt"`
		} `json:"versionInfo"`
	} `json:"app"`
}