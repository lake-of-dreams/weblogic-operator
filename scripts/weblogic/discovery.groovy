// validate the bootstrap model file


import oracle.fmwplatform.envspec.environment.OfflineEnvironment;
import oracle.fmwplatform.envspec.environment.topology.DiscoveryOptions;
import oracle.fmwplatform.envspec.model.EnvironmentModel;
import oracle.fmwplatform.envspec.model.EnvironmentModelBuilder;
import oracle.fmwplatform.envspec.model.targets.ModelTarget;
import oracle.fmwplatform.envspec.model.targets.ModelTargetFactory;

def builder = new EnvironmentModelBuilder("/u01/oracle")
def bootStrapModel = builder.buildFromOfflineDomainBootstrapParams("firstdomain", "/u01/oracle", "/u01/oracle/user_projects/domains/firstdomain");

def targets = new ArrayList<ModelTarget>();
targets.add(ModelTargetFactory.createDomainTarget("firstdomain"));

def options = new DiscoveryOptions(false, false, OfflineEnvironment.class);
def model = builder.populateFromEnvironment(bootStrapModel, options, targets);

print model.getTopology()
