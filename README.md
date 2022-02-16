# GO sample rest application

## Features, concepts: 
- [ ] authentication
- [x] CRUD + DB connection (mongodb)
- [x] HTTP routing (mux)
- [x] logging
- [x] CSV file upload 
- [x] goroutines 
- [x] channels

## How to use
### Clone repository

```
git clone https://github.com/mj-hagonoy/mjh-go-rest.git
cd mjh-go-rest
```


### Update config file
Refer to config.yaml
```
host: ""
port: 8080
database: 
  host: "0.0.0.0"
  port: 27017
  username: ""
  password: ""
  dbname: "<db_name>"
  sslmode: "disable"
directory:
  import_users: "<upload directory for csv files - import users>"
  mail_templates: <email templates directory>"
mail:
  from: "...@gmail.com"
  smtp_host: ""
  smtp_port: ""
  smtp_user : ""
  smtp_pwd: ""
log:
  log_dir: "<log directory>"
```


### Build, run locally
```
make build
make local
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
