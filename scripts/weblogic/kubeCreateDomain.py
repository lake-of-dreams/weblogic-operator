import sys, socket

def addCluster(clusterName):
    cd('/')
    clusterId = create(clusterName, 'Cluster')
    cd('/')
    return clusterId;


def addMachine(machineName):
    cd('/')
    machine = create(machineName, 'Machine')
    cd('/')
    return machine;


def addNodeManager(nmHost):
    cd('/')
    cd('/Machines/' + nmHost)
    nm = create(nmHost, 'NodeManager')
    nm.setListenAddress(nmHost)
    nm.setListenPort(nodeMgrPort)
    nm.setDebugEnabled(true)
    cd('/')
    print(nm)
    return;


def addManagedServer(serverName, serverPort, serverHost, serverMac):
    cd('/')
    create(serverName, 'Server')
    cd('Servers/' + serverName)
    set('ListenPort', serverPort)
    set('ListenAddress', serverHost)
    set('Machine', serverMac)
    if clusterExist:
        set('Cluster', cluster)

    cd('/')
    return;


def addAdminServer():
    cd('/')
    set('Name', domainName)
    cd('/Servers/AdminServer')
    set('Name', 'AdminServer')
    set('ListenPort', adminPort)
    set('ListenAddress', adminHost)

    cd('/')

    cd('Security/' + domainName + '/User/weblogic')
    cmo.setName(username)
    cmo.setPassword(password)

    cd('/')
    macA = create(adminHost, 'Machine')
    cd('/Machines/' + adminHost)
    nm = create(adminHost, 'NodeManager')
    nm.setListenAddress(adminHost)
    nm.setListenPort(nodeMgrPort)
    nm.setDebugEnabled(true)

    cd('/')
    cd('/Servers/AdminServer')
    set('Machine', macA)
    return;

### MAIN

try:
    clusterExist = False

    domainDir = sys.argv[1]
    username = sys.argv[2]
    password = sys.argv[3]

    adminHost = sys.argv[4]
    adminPort = int(sys.argv[5])

    managedHostCount = int(sys.argv[6])

    nodeMgrPort = int(sys.argv[7])

    if domainDir.endsWith("/"):
        domainDir[-1].replace("/", "")

    data = domainDir.split("/");
    domainName = data(len(data) - 1)

    print "DOMAIN NAME: " + domainName
    print "DOMAIN HOME: " + domainDir
    print "USERNAME: " + username
    print "PASSWORD: " + password

    print "ADMIN SERVER HOST: " + adminHost
    print "ADMIN SERVER PORT: " + str(adminPort)
    print "NODE MANAGER PORT: " + str(nodeMgrPort)

    selectTemplate('Basic WebLogic Server Domain')
    loadTemplates()

    addAdminServer()

    if managedHostCount > 1:
        clusterExist = True

    if clusterExist:
        cluster = addCluster('Cluster1')

    initport = adminPort;
    list = [];
    for x in range(1, managedHostCount + 1):
        initport += 2
        servername = 'Server' + x
        host = 'localhost'
        dictServer = {'ServerName': servername, 'Port': initport, 'Host': host}
        list.append(dict)

        mac = addMachine(host)
        addNodeManager(host)
        addManagedServer(servername, initport, host, mac)

    writeDomain(domainDir)
    closeTemplate()
    print
    "Domain " + domainName + " Created Successfully!"

except Exception, e:
    e.printStackTrace()
    dumpStack()
    raise("Create Domain Failed")
