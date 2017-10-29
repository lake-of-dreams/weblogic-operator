import os
import sys

def addNodeManager(machineName):
    cd('/')
    machine = create(machineName, 'Machine')

    cd('/Machines/' + machineName)
    nm = create(machineName, 'NodeManager')
    nm.setListenAddress('localhost')
    nm.setListenPort(8886)
    nm.setDebugEnabled(true)
    cd('/')
    print(nm)
    return machine;


def addManagedServer(serverName, serverPort):
    mac=addNodeManager('Machine-'+serverName)

    cd('/')
    create(serverName, 'Server')
    cd('Servers/' + serverName)
    set('ListenPort', serverPort)
    set('ListenAddress', '')
    set('Machine', mac)

    cd('/')
    return;


### MAIN

try:
    # Variable Definitions
    # ======================
    myPodName = sys.argv[1]
    oracleHome = sys.argv[2]
    domainName = sys.argv[3]
    domainHome = sys.argv[4]
    managedServerCount = int(sys.argv[5])
    adminPort = int(sys.argv[6])
    username = sys.argv[7]
    password = sys.argv[8]

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

    cd('/')
    macA = create('Machine-AdminServer', 'Machine')
    cd('/Machines/' + 'Machine-AdminServer')
    nm = create('Machine-AdminServer', 'NodeManager')
    nm.setListenAddress(myPodName)
    nm.setListenPort(8886)
    nm.setDebugEnabled(true)

    cd('/')
    cd('/Servers/AdminServer')
    set('Machine', macA)

    # Create Managed Servers
    # =====================================
    port = adminPort;
    serverlist = [];
    dictServer = {"serverName": "AdminServer", "port": adminPort, "host": "localhost", "podName": myPodName}
    serverlist.append(dictServer)
    for x in range(1, managedServerCount + 1):
        port += 2
        servername = 'managedserver-' + str((x - 1))
        host = 'localhost'
        dictServer = {"serverName": servername, "port": port, "host": host, "podName": ""}
        serverlist.append(dictServer)

        addManagedServer(servername, port)

    # Write Domain
    # ============
    writeDomain(domainHome)
    closeTemplate()
    print "Domain Created Successfully "

    # Save Server List
    # ================
    serverListFile = '%s/serverList.json' % domainHome
    os.system("touch %s" % serverListFile)
    file = open(serverListFile, "w")
    file.write("[")
    for item in serverlist:
        file.write("\n%s," % item)
    file.write("]\n")
    file.close()

    file = open(serverListFile, "r")
    filedata = file.read()
    filedata = filedata.replace('\'', '\"')
    filedata = filedata.replace(',]', '\n]')
    file = open(serverListFile, "w")
    file.write(filedata)

    # Exit WLST
    # =========
    exit()

except Exception, e:
    e.printStackTrace()
    dumpStack()
    raise ("Create Domain Failed")
