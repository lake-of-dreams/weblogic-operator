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
    clusterExist = False
    oracleHome = sys.argv[1]
    domainName = sys.argv[2]
    domainHome = sys.argv[3]
    managedServerCount = int(sys.argv[4])
    adminPort = int(sys.argv[5])
    username = sys.argv[6]
    password = sys.argv[7]

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
        servername = 'managedserver-' + str((x - 1))
        host = 'localhost'
        dictServer = {'ServerName': servername, 'Port': port, 'Host': host}
        serverlist.append(dictServer)

        addManagedServer(servername, port)

    # Write Domain
    # ============
    writeDomain(domainHome)
    closeTemplate()
    print "Domain Created Successfully "

    # Save Server List
    # ================
    serverListFile = '%s/serverList.txt' % domainHome
    os.system("touch %s" % serverListFile)
    file = open(serverListFile, "w+")
    for item in serverlist:
        file.write("%s\n" % item)
    file.close()

    # Exit WLST
    # =========
    exit()

except Exception, e:
    e.printStackTrace()
    dumpStack()
    raise ("Create Domain Failed")
