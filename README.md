This program replaces docker images with images from the last successful Jenkins build.

Before first run:
- Create a `config.json` file using the `example.config.json` file and fill all <...> places.
- The files `config.json` and `main.go` must be located in same directory.

Run:
- In the terminal go to the directory with the `main.go` file.
- Execute: `go run main.go`

Example `config.json`:
```json
{
    "credential": {
        "username": "Markiz",
        "token": "8721af09df9cb09b8beee7987b98717357"
    },
    "composePath": "D:/dockerFiles/docker-compose.yaml",
    "trackedJobs": [
        {
            "url": "http://myjenkinsserver.com/job/backend/job/master/"
        },
        {
            "url": "http://myjenkinsserver.com/job/frontend/job/master/"
        }
    ]
}
```
Example `docker-compose.yaml`:
```yaml
networks:
  local:
    driver: bridge
    ipam:
      config:
        - gateway: 123.45.0.1
          subnet: 123.45.0.0/16
      driver: default

services:
  backend:
    image: nexus.myserver.com/backend:2.0.3
    environment:
      TERM: xterm
      DEBUG: "true"
      DEBUG_PORT: "8787"
    networks:
      local:
        ipv4_address: 123.45.0.2
    ports:
      - 127.0.0.1:8080:8080
      - 127.0.0.1:8787:8787
    volumes:
      - "D:/dockerFiles/backend-data:/opt/jboss/server/default:rw"

  frontend:
    image: nexus.myserver.com/frontend:5.1.7
    environment:
      TERM: xterm
    depends_on:
      - backend
    networks:
      local:
        ipv4_address: 123.45.0.3
    ports:
      - 127.0.0.1:8080:8080

version: "2.0"
```

Do request GET `config.TrackedJobs[0].URL` and parse response:
```json
{
    "_class": "org.jenkinsci.plugins.workflow.job.WorkflowJob",
    "actions": [
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {
            "_class": "hudson.plugins.tasks.TasksProjectAction"
        },
        {},
        {},
        {},
        {
            "_class": "com.cloudbees.plugins.credentials.ViewCredentialsAction"
        }
    ],
    "description": null,
    "displayName": "ink",
    "displayNameOrNull": null,
    "fullDisplayName": "backend » master",
    "fullName": "backend/master",
    "name": "ink",
    "url": "http://myjenkinsserver.com/job/backend/job/master/",
    "buildable": true,
    "builds": [
        {
            "_class": "org.jenkinsci.plugins.workflow.job.WorkflowRun",
            "number": 157,
            "url": "http://myjenkinsserver.com/job/backend/job/master/157/"
        },
        {
            "_class": "org.jenkinsci.plugins.workflow.job.WorkflowRun",
            "number": 156,
            "url": "http://myjenkinsserver.com/job/backend/job/master/156/"
        },
        {
            "_class": "org.jenkinsci.plugins.workflow.job.WorkflowRun",
            "number": 133,
            "url": "http://myjenkinsserver.com/job/backend/job/master/133/"
        }
    ],
    "color": "blue",
    "firstBuild": {
        "_class": "org.jenkinsci.plugins.workflow.job.WorkflowRun",
        "number": 133,
        "url": "http://myjenkinsserver.com/job/backend/job/master/133/"
    },
    "healthReport": [
        {
            "description": "Стабильность сборок: Среди последних сборок провалившихся нет.",
            "iconClassName": "icon-health-80plus",
            "iconUrl": "health-80plus.png",
            "score": 100
        }
    ],
    "inQueue": false,
    "keepDependencies": false,
    "lastBuild": {
        "_class": "org.jenkinsci.plugins.workflow.job.WorkflowRun",
        "number": 157,
        "url": "http://myjenkinsserver.com/job/backend/job/master/157/"
    },
    "lastCompletedBuild": {
        "_class": "org.jenkinsci.plugins.workflow.job.WorkflowRun",
        "number": 157,
        "url": "http://myjenkinsserver.com/job/backend/job/master/157/"
    },
    "lastFailedBuild": null,
    "lastStableBuild": {
        "_class": "org.jenkinsci.plugins.workflow.job.WorkflowRun",
        "number": 157,
        "url": "http://myjenkinsserver.com/job/backend/job/master/157/"
    },
    "lastSuccessfulBuild": {
        "_class": "org.jenkinsci.plugins.workflow.job.WorkflowRun",
        "number": 157,
        "url": "http://myjenkinsserver.com/job/backend/job/master/157/"
    },
    "lastUnstableBuild": null,
    "lastUnsuccessfulBuild": null,
    "nextBuildNumber": 158,
    "property": [
        {
            "_class": "jenkins.model.BuildDiscarderProperty"
        },
        {
            "_class": "org.jenkinsci.plugins.workflow.job.properties.DisableConcurrentBuildsJobProperty"
        },
        {
            "_class": "org.jenkinsci.plugins.workflow.multibranch.BranchJobProperty",
            "branch": {}
        }
    ],
    "queueItem": null,
    "concurrentBuild": false,
    "resumeBlocked": false
}
```
Then we get `response.lastSuccessfulBuild.number` and `response.lastSuccessfulBuild.url`

Do request GET `response.lastSuccessfulBuild.url` and parse response:
```json
{
    "_class": "org.jenkinsci.plugins.workflow.job.WorkflowRun",
    "actions": [
        {
            "_class": "hudson.model.CauseAction",
            "causes": [
                {
                    "_class": "hudson.model.Cause$UserIdCause",
                    "shortDescription": "Started by user Markiz",
                    "userId": "Markiz",
                    "userName": "Markiz"
                }
            ]
        },
        {
            "_class": "jenkins.metrics.impl.TimeInQueueAction",
            "blockedDurationMillis": 0,
            "blockedTimeMillis": 0,
            "buildableDurationMillis": 0,
            "buildableTimeMillis": 4,
            "buildingDurationMillis": 99745,
            "executingTimeMillis": 98221,
            "executorUtilization": 0.98,
            "subTaskCount": 1,
            "waitingDurationMillis": 0,
            "waitingTimeMillis": 0
        },
        {
            "_class": "jenkins.scm.api.SCMRevisionAction"
        },
        {},
        {},
        {
            "_class": "hudson.plugins.git.util.BuildData",
            "buildsByBranchName": {
                "ink": {
                    "_class": "hudson.plugins.git.util.Build",
                    "buildNumber": 157,
                    "buildResult": null,
                    "marked": {
                        "SHA1": "c3d0a66321c46b40cef5d07a8526987e80086ad7",
                        "branch": [
                            {
                                "SHA1": "c3d0a66321c46b40cef5d07a8526987e80086ad7",
                                "name": "master"
                            }
                        ]
                    },
                    "revision": {
                        "SHA1": "c3d0a66321c46b40cef5d07a8526987e80086ad7",
                        "branch": [
                            {
                                "SHA1": "c3d0a66321c46b40cef5d07a8526987e80086ad7",
                                "name": "master"
                            }
                        ]
                    }
                }
            },
            "lastBuiltRevision": {
                "SHA1": "c3d0a66321c46b40cef5d07a8526987e80086ad7",
                "branch": [
                    {
                        "SHA1": "c3d0a66321c46b40cef5d07a8526987e80086ad7",
                        "name": "master"
                    }
                ]
            },
            "remoteUrls": [
                "git@lab:coolapp/backend.git"
            ],
            "scmName": ""
        },
        {
            "_class": "hudson.plugins.git.GitTagAction"
        },
        {},
        {},
        {
            "_class": "org.jenkinsci.plugins.workflow.cps.EnvActionImpl"
        },
        {},
        {},
        {
            "_class": "hudson.plugins.tasks.TasksResultAction"
        },
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {
            "_class": "org.jenkinsci.plugins.pipeline.modeldefinition.actions.RestartDeclarativePipelineAction"
        },
        {},
        {
            "_class": "org.jenkinsci.plugins.workflow.job.views.FlowGraphAction"
        },
        {},
        {},
        {}
    ],
    "artifacts": [
        {
            "displayPath": null,
            "fileName": "backend-2.0.3.pom",
            "relativePath": "com/mycompany/coolapp/master/backend/2.0.3/backend-2.0.3.pom"
        }
    ],
    "building": false,
    "description": "nexus.myserver.com/backend:2.0.4\n",
    "displayName": "#157",
    "duration": 99745,
    "estimatedDuration": 101983,
    "executor": null,
    "fullDisplayName": "backend » backend #157",
    "id": "157",
    "keepLog": false,
    "number": 157,
    "queueId": 7050457,
    "result": "SUCCESS",
    "timestamp": 1602772390245,
    "url": "http://myjenkinsserver.com/job/backend/job/master/157/",
    "changeSets": [],
    "culprits": [],
    "nextBuild": null,
    "previousBuild": {
        "number": 156,
        "url": "http://myjenkinsserver.com/job/backend/job/master/156/"
    }
}
```
And we get `response.description` that contains docker images names.

Next we find the string `nexus.myserver.com/backend:2.0.3` using the regexp string `nexus.myserver.com/backend:\S` in `docker-compose.yaml` and replace it with `nexus.myserver.com/backend:2.0.4` from `response.description`.
