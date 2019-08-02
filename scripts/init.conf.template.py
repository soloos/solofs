#!/usr/bin/env python3
# -*- coding: utf-8 -*-

temp = '''
{{
    "DataNodeSRPCPeerID":"DATA_NODE_SRPC_{Index}_00000000000000000000000000000000000000000000000",
    "SRPCListenAddr":"192.168.56.100:{PortPrefix}01",
    "SRPCServeAddr":"192.168.56.100:{PortPrefix}01",
    "DataNodeWebPeerID":"DATA_NODE_Web_{Index}_000000000000000000000000000000000000000000000000",
    "WebServer": {{
        "ServeStr":"http://192.168.56.100:{PortPrefix}02",
        "ListenStr":"192.168.56.100:{PortPrefix}02"
    }},
    "DataNodeLocalFSRoot":"/tmp/sdfs_test.data.0{Index}",
    "PProfListenAddr":"192.168.56.100:{PortPrefix}03"
}}
'''
for i in range(4):
    filepath = './scripts/conf_template/datanode.{Index}.json'.format(Index=i)
    with open(filepath, 'w+') as f:
        content = temp.format(Index=i, PortPrefix=106+i)
        f.write(content)

