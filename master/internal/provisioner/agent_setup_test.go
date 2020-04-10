package provisioner

import (
	"encoding/base64"
	"testing"

	"gotest.tools/assert"

	"github.com/determined-ai/determined/master/pkg/etc"
)

func TestAgentSetupScript(t *testing.T) {
	err := etc.SetRootPath("../../static/srv/")
	assert.NilError(t, err)

	encoded := base64.StdEncoding.EncodeToString([]byte("sleep 5\n echo \"hello world\""))
	conf := agentSetupScriptConfig{
		MasterHost:          "test.master",
		MasterPort:          "8080",
		StartupScriptBase64: encoded,
		AgentDockerImage:    "test_docker_image",
		AgentDockerRuntime:  "nvidia",
		AgentNetwork:        "default",
		AgentID:             "test.id",
	}
	expected := `#!/bin/bash

mkdir -p /usr/local/determined
echo c2xlZXAgNQogZWNobyAiaGVsbG8gd29ybGQi | base64 --decode > /usr/local/determined/startup_script
echo "#### PRINTING STARTUP SCRIPT START ####"
cat /usr/local/determined/startup_script
echo "#### PRINTING STARTUP SCRIPT END ####"
chmod +x /usr/local/determined/startup_script
/usr/local/determined/startup_script

docker run --init --name determined-agent  --restart always --network default --runtime=nvidia \
    -e DET_AGENT_ID="test.id" \
    -e DET_MASTER_HOST="test.master" \
    -e DET_MASTER_PORT="8080" \
    -v /var/run/docker.sock:/var/run/docker.sock \
    "test_docker_image"
`

	res := string(mustMakeAgentSetupScript(conf))
	assert.Equal(t, res, expected)
}
