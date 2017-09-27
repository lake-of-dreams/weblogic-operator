import os, sys

def addCluster(clusterName):
    cd('/')
    clusterId = create(clusterName, 'Cluster')
    cd('/')
    return clusterId;


def addManagedServer(serverName, serverPort):
    cd('/')
    create(serverName, 'Server')
    cd('Servers/' + serverName)
    set('ListenPort', serverPort)
    set('ListenAddress', '')
    if clusterExist:
        set('Cluster', cluster)

    cd('/')
    return;


### MAIN

try:
    # Variable Definitions
    # ======================
    oracleHome = os.environ.get("ORACLE_HOME", "/u01/oracle")
    domainName = os.environ.get("DOMAIN_NAME", "basedomain")
    domainHome = os.environ.get("DOMAIN_HOME", '/u01/oracle/user_projects/domains/%s' % domainName)
    managedServerCount = int(os.environ.get("MANAGED_SERVER_COUNT", "1"))
    adminPort = int("7001")
    username = 'weblogic'
    password = 'welcome1'

    print('ORACLE_HOME              : [%s]' % oracleHome);
    print('DOMAIN_NAME              : [%s]' % domainName);
    print('DOMAIN_HOME              : [%s]' % domainHome);
    print('MANAGED_SERVER_COUNT     : [%s]' % managedServerCount);
    print('ADMIN_PORT               : [%s]' % adminPort);
    print('USERNAME                 : [%s]' % username);
    print('PASSWORD                 : [%s]' % password);

    # Open default domain template
    # ======================
    selectTemplate('Basic WebLogic Server Domain')
    loadTemplates()

    set('Name', domainName)
    setOption('DomainName', domainName)

    # Configure the Administration Server and SSL port.
    # =========================================================
    cd('/Servers/AdminServer')
    set('ListenAddress', '')
    set('ListenPort', adminPort)

    # Define the user password for weblogic
    # =====================================
    cd('/')

    cd('Security/' + domainName + '/User/weblogic')
    cmo.setName(username)
    cmo.setPassword(password)

    # Create Managed Servers
    # =====================================
    if managedServerCount > 1:
        clusterExist = True

    if clusterExist:
        cluster = addCluster('cluster-0')

    port = adminPort;
    serverlist = [];
    for x in range(1, managedServerCount + 1):
        port += 2
        servername = 'managedserver-' + (x - 1)
        host = 'localhost'
        dictServer = {'ServerName': servername, 'Port': port, 'Host': host}
        serverlist.append(dictServer)

        addManagedServer(servername, port)

    serverListFile = '%s/serverList.txt' % domainHome
    file = open(serverListFile, 'w')
    for item in serverlist:
        file.write("%s\n" % item)
    file.close()

    # Write Domain
    # ============
    writeDomain(domainHome)
    closeTemplate()
    print "Domain Created Successfully "

    # Exit WLST
    # =========
    exit()

except Exception, e:
    e.printStackTrace()
    dumpStack()
    raise ("Create Domain Failed")
