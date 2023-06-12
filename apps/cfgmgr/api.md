Each application (i.e. every concurrent lambda), loads the configuraions during cold start (i.e. before handling any reuqest). This data is cached for the lifetime of the application and used for all request handling. 

#### Servers:
    Local:  http://localhost:5001/v1
    UAT:    http://cfgmgr.uat.myorg.com/v2

