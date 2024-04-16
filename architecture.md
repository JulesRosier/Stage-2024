# Graph
```mermaid
graph TD
 subgraph "Main (Start)"
 A("Repeating")
 end
 subgraph B["Gent Open Data"]
 C("Stations API Endpoint")
 end
 subgraph TR["Transaction"]
 subgraph D["Create/Update"]
 DECI(" ")
 CREATE("Create")
 UPD("Update")
 G("ColumnChange")
 CHANGEDET("ChangeDetected")
 end
 subgraph PG["Postgres"]
 HSD("Historical Station Data")
 OLTP("OLTP")
 OUT("Outbox")
  end
 end
 subgraph L["Listener"]
 M("Fetch")
 N("Produce")
 end
 RP("Redpanda")

 subgraph Q["Bike Event"]
 GEN("Generator")
 end
 subgraph S["Webserver"]
 T("Dashboard")
 end

 A--> A
 A --> B
 B ---|JSON Data| DECI
 UPD --> G
 CREATE --> PG
 UPD -.-> OLTP
 G --> CHANGEDET
 CHANGEDET --> PG
 OUT <-.->|30s| L
 M -->|Queue| N
 N --> RP
 GEN <-.->|Fetches data every 60m| PG
 RP <-.-> S
 GEN -.-> PG
 DECI-->CREATE
 DECI-->UPD
 ```