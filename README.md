# GO sample rest application

## Features, concepts: 
- [ ] authentication
- [x] CRUD + DB connection (mongodb)
- [x] HTTP routing (mux)
- [x] logging
- [x] goroutines 
- [x] channels

### Google Cloud services used
- [x] Google Cloud Storage
- [x] Google Cloud PubSub

## How to use
### Clone repository

```
git clone https://github.com/mj-hagonoy/mjh-go-rest.git
cd mjh-go-rest
```


### Update config file
Refer to [config.yaml](config.yaml)
```
host: "localhost"
port: 8080
database: 
  host: "0.0.0.0"
  port: 27017
  username: ""
  password: ""
  dbname: "mjh"
  sslmode: "disable"
directory:
  import_users: "<upload directory for csv files - import users>"
  mail_templates: <email templates directory>"
mail:
  from: "...@gmail.com"
  smtp_host: "smtp.mailtrap.io"
  smtp_port: "2525"
  smtp_user : "<user>"
  smtp_pwd: "<pwd>"
log:
  log_dir: "<log directory>"
file_storage:
  default: "google_cloud"
  google_cloud:
    project_id: "<gcp project ID>"
    bucket_name: "<gcp bucket name>"
    upload_path: "<upload path = [empty or subfolder]>"
credentials:
  google_app_creds : "<file path to GOOGLE_APPLICATION_CREDENTIALS>"
messaging:
  google_cloud: 
    project_id: "<gcp project ID>"
    topic_id: "<gcp bucket name>"
```

### Build
```
go build ./main.go ./jobworker.go
```

### Run web service
```
./main --config config.yaml --type web
```

### Run job service
```
./main --config config.yaml --type job
```

## APIs
| API | Method | Request Body | Response | Description |
| --- | ------ | ------------ | -------- | ----------- |
| /api/v1/users/import | POST | `file` : csv file | (json) batch job info [id, status] | accepts csv file with users info [first name, last name, email] |
| /api/v1/users | GET | | (json) users | returns a list of users |
| /api/v1/jobs/:id | GET | |(json) job info | returns information of the job indicated by `id` |

## TODO
- [ ] docker container
- [ ] db initialization scripts
