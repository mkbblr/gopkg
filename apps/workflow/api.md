Workflow is a server application that implements business workflows. This is a monolithic application i.e. all the API, events, tools are built into single binary and can be invoked using different command line arguments, API integrations etc. 

Even though the it is monolithic application, the design is such that we can deploy it as a distributed application with each service integrated to different set of endpoints. 

It is inspired by the blog post [Designing A Workflow Engine](https://exceptionnotfound.net/designing-a-workflow-engine-database-part-1-introduction-and-purpose/).


#### Servers:
    Local:  http://localhost:5001/v1
    UAT:    http://workflow.uat.myorg.com/v2

