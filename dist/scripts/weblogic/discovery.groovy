import oracle.fmwplatform.envspec.environment.OfflineEnvironment;
import oracle.fmwplatform.envspec.environment.OnlineEnvironment;
import oracle.fmwplatform.envspec.environment.topology.DiscoveryOptions;
import oracle.fmwplatform.envspec.model.EnvironmentModel;
import oracle.fmwplatform.envspec.model.EnvironmentModelBuilder;
import oracle.fmwplatform.envspec.model.targets.ModelTarget;
import oracle.fmwplatform.envspec.model.targets.ModelTargetFactory;
import oracle.fmwplatform.credentials.credential.Credentials;
import oracle.fmwplatform.credentials.wallet.WalletStoreProvider;

def credentials = new Credentials();
credentials.setCredential("WLS/ADMIN", "weblogic", "welcome1".toCharArray());

def walletStoreProvider = new WalletStoreProvider("/u01/oracle/user_projects/firstdomainWallet", "welcome1".toCharArray());
walletStoreProvider.createWallet();
walletStoreProvider.storeCredentials(credentials);
walletStoreProvider.closeWallet(false);

def builder = new EnvironmentModelBuilder("/u01/oracle")
//def bootStrapModel = builder.buildFromOfflineDomainBootstrapParams("firstdomain", "/u01/oracle", "/u01/oracle/user_projects/domains/firstdomain");
def bootStrapModel = builder.buildFromOnlineDomainBootstrapParams("firstdomain", "/u01/oracle", "t3://firstdomain-dw86k:7001", credentials.getCredential("WLS/ADMIN"));

def targets = new ArrayList<ModelTarget>();
targets.add(ModelTargetFactory.createDomainTarget("firstdomain"));

//def options = new DiscoveryOptions(false, false, OfflineEnvironment.class);
def options = new DiscoveryOptions(false, true, OnlineEnvironment.class);
def model = builder.populateFromEnvironment(bootStrapModel, options, targets);

print model.getTopology()
