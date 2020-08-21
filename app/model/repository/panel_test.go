package repository

import (
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"log"
	"os"
	"testing"
	"time"
	"voyageone.com/dp/infrastructure/model/global"
)

func init() {
	cluster := gocql.NewCluster("127.0.0.1:9042")
	cluster.Consistency = gocql.Quorum
	cluster.NumConns = 3
	var err error
	global.CqlSession, err = gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		log.Fatal("CQL Session Create Failed!", err)
	}
	global.DPLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Llongfile)
}

func TestCreateOrUpdateStatusHistory(t *testing.T) {
	inputs := []StatusHistory{
		{
			AppName:         "app_a",
			StatusChangedTo: "deploying",
			Time:            time.Now(),
		},
		{
			AppName:         "app_b",
			StatusChangedTo: "starting",
			Time:            time.Now(),
		},
		{
			AppName:         "app_c",
			StatusChangedTo: "stopping",
			Time:            time.Now(),
		},
	}
	for _, i := range inputs {
		err := i.CreateOrUpdate()
		if err != nil {
			t.Error(err)
		}
	}
}

func TestGetStoppedAppsList(t *testing.T) {
	jobs, err := GetStoppedAppsList()
	if err != nil {
		t.Error(err)
	} else {
		println(jobs)
	}
}

func TestStoppedAppCache_CreateOrUpdate(t *testing.T) {
	c := StoppedAppCache{
		AppName:      "openvms-restapi-dp",
		NomadJobJson: `{"Stop":false,"Region":"vo-local","Namespace":"default","ID":"openvms-restapi-dp","ParentID":"","Name":"openvms-restapi-dp","Type":"service","Priority":50,"AllAtOnce":false,"Datacenters":["host70","host80"],"Constraints":null,"Affinities":null,"TaskGroups":[{"Name":"group","Count":1,"Constraints":[{"LTarget":"","RTarget":"true","Operand":"distinct_hosts"},{"LTarget":"${meta.for}","RTarget":"gateway","Operand":"!="}],"Affinities":[{"LTarget":"${node.datacenter}","RTarget":"host70","Operand":"=","Weight":90}],"Tasks":[{"Name":"task","Driver":"java","User":"","Lifecycle":null,"Config":{"jar_path":"local/openvms-restapi-dp.jar","jvm_options":["-javaagent:local/agent/skywalking-agent.jar","-Xmx512m","-Xms512m","-Dlogging.config=local/logback.xml","-agentlib:jdwp=transport=dt_socket,server=y,suspend=n,address=${NOMAD_PORT_debug}","-Dserver.port=${NOMAD_PORT_http}"]},"Constraints":null,"Affinities":null,"Env":null,"Services":[{"Id":"","Name":"openvms-restapi-dp","Tags":["green"],"CanaryTags":["canary"],"EnableTagOverride":false,"PortLabel":"http","AddressMode":"auto","Checks":[{"Id":"","Name":"Spring Actuator Health","Type":"http","Command":"","Args":null,"Path":"/health","Protocol":"","PortLabel":"http","Expose":false,"AddressMode":"","Interval":12000000000,"Timeout":4000000000,"InitialStatus":"","TLSSkipVerify":false,"Header":null,"Method":"","CheckRestart":null,"GRPCService":"","GRPCUseTLS":false,"TaskName":""}],"CheckRestart":null,"Connect":null,"Meta":null,"CanaryMeta":null,"TaskName":""}],"Resources":{"CPU":500,"MemoryMB":1500,"DiskMB":0,"Networks":[{"Mode":"","Device":"","CIDR":"","IP":"","MBits":1000,"DNS":null,"ReservedPorts":null,"DynamicPorts":[{"Label":"http","Value":0,"To":0,"HostNetwork":""},{"Label":"debug","Value":0,"To":0,"HostNetwork":""}]}],"Devices":null,"IOPS":0},"RestartPolicy":{"Interval":1800000000000,"Attempts":1,"Delay":15000000000,"Mode":"fail"},"Meta":null,"KillTimeout":5000000000,"LogConfig":{"MaxFiles":10,"MaxFileSizeMB":10},"Artifacts":[{"GetterSource":"http://10.0.0.70:4567/voerp-dev/openvms-restapi-dp/openvms-restapi-dp-v202008101522.jar","GetterOptions":null,"GetterMode":"file","RelativeDest":"local/openvms-restapi-dp.jar"},{"GetterSource":"http://10.0.0.70:4567/skywalking/skwalking-agent-6.6.0-staging.zip","GetterOptions":null,"GetterMode":"any","RelativeDest":"local/"},{"GetterSource":"http://10.0.0.70:4567/skywalking/skywalking-agent.config.tpl","GetterOptions":null,"GetterMode":"file","RelativeDest":"local/agent.config.tpl"},{"GetterSource":"http://10.0.0.70:4567/logback/logback-nomad-skywalking.xml","GetterOptions":null,"GetterMode":"file","RelativeDest":"local/logback.xml"}],"Vault":null,"Templates":[{"SourcePath":"local/agent.config.tpl","DestPath":"local/agent/config/agent.config","EmbeddedTmpl":"","ChangeMode":"noop","ChangeSignal":"","Splay":5000000000,"Perms":"0644","LeftDelim":"{{","RightDelim":"}}","Envvars":false,"VaultGrace":0},{"SourcePath":"","DestPath":"/tmp/restart_flag","EmbeddedTmpl":"{{ key \"services/openvms-restapi-dp/restart\" }}","ChangeMode":"restart","ChangeSignal":"","Splay":30000000000,"Perms":"0644","LeftDelim":"{{","RightDelim":"}}","Envvars":false,"VaultGrace":0}],"DispatchPayload":null,"VolumeMounts":[{"Volume":"logs","Destination":"/data/logs","ReadOnly":false,"PropagationMode":"private"},{"Volume":"share","Destination":"/data/share","ReadOnly":false,"PropagationMode":"private"}],"Leader":false,"ShutdownDelay":0,"KillSignal":"","Kind":""}],"Spreads":null,"Volumes":{"logs":{"Name":"logs","Type":"host","Source":"logs","ReadOnly":false,"MountOptions":null},"share":{"Name":"share","Type":"host","Source":"share","ReadOnly":false,"MountOptions":null}},"RestartPolicy":{"Interval":1800000000000,"Attempts":1,"Delay":15000000000,"Mode":"fail"},"ReschedulePolicy":{"Attempts":0,"Interval":0,"Delay":30000000000,"DelayFunction":"exponential","MaxDelay":3600000000000,"Unlimited":true},"EphemeralDisk":{"Sticky":false,"Migrate":false,"SizeMB":300},"Update":{"Stagger":30000000000,"MaxParallel":1,"HealthCheck":"checks","MinHealthyTime":30000000000,"HealthyDeadline":300000000000,"ProgressDeadline":600000000000,"Canary":0,"AutoRevert":false,"AutoPromote":false},"Migrate":{"MaxParallel":1,"HealthCheck":"checks","MinHealthyTime":30000000000,"HealthyDeadline":300000000000},"Networks":null,"Meta":null,"Services":null,"ShutdownDelay":null,"StopAfterClientDisconnect":null,"Scaling":null}],"Update":{"Stagger":30000000000,"MaxParallel":1,"HealthCheck":"","MinHealthyTime":0,"HealthyDeadline":0,"ProgressDeadline":0,"Canary":0,"AutoRevert":false,"AutoPromote":false},"Multiregion":null,"Spreads":null,"Periodic":null,"ParameterizedJob":null,"Dispatched":false,"Payload":null,"Reschedule":null,"Migrate":null,"Meta":{"last_manual_restart":"0"},"ConsulToken":"","VaultToken":"","VaultNamespace":null,"NomadTokenID":null,"Status":"running","StatusDescription":"","Stable":true,"Version":15,"SubmitTime":1597127439908217740,"CreateIndex":477412,"ModifyIndex":479558,"JobModifyIndex":479547}`,
	}
	err := c.CreateOrUpdate()
	if err != nil {
		t.Error(err)
	}
}
