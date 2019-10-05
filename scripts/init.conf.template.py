#!/usr/bin/env python3
# -*- coding: utf-8 -*-

temp = '''
{{
    "SolodnSrpcPeerID":"DATA_NODE_SRPC_{Index}_00000000000000000000000000000000000000000000000",
    "SrpcListenAddr":"192.168.56.100:{PortPrefix}01",
    "SrpcServeAddr":"192.168.56.100:{PortPrefix}01",
    "SolodnWebPeerID":"DATA_NODE_Web_{Index}_000000000000000000000000000000000000000000000000",
    "WebServer": {{
        "ServeStr":"http://192.168.56.100:{PortPrefix}02",
        "ListenStr":"192.168.56.100:{PortPrefix}02"
    }},
    "SolodnLocalFsRoot":"/tmp/soloos_test.data.0{Index}",
    "PProfListenAddr":"192.168.56.100:{PortPrefix}03"
}}
'''
for i in range(4):
    filepath = './scripts/conf_template/solodn.{Index}.json'.format(Index=i)
    with open(filepath, 'w+') as f:
        content = temp.format(Index=i, PortPrefix=106+i)
        f.write(content)

