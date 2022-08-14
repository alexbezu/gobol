```mermaid
graph LR;
 client([client])-. 23:23567 .->ingress[3270 server];
 ingress-->|3322:3322|service[immudb];
 subgraph cluster[IMS]
 ingress;
 service<-->|3322:3322|pod1[Transaction1];
 service-->|3322:3322|pod2[Transaction2];
 service-->pod4(StateFulSet);
 ingress<-->pod3(PersistVolume);
 end
 classDef plain fill:#ddd,stroke:#fff,stroke-width:4px,color:#000;
 classDef k8s fill:#326ce5,stroke:#fff,stroke-width:4px,color:#fff;
 classDef cluster fill:#fff,stroke:#bbb,stroke-width:2px,color:#326ce5;
 class ingress,service,pod1,pod2 k8s;
 class client plain;
 class cluster cluster;
 ```