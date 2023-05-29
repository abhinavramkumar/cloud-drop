# Cloud Drop

Simple script to upload any number of file(s) to Amazon S3

### Clone

```bash
git clone https://github.com/abhinav-ramkumar/cloud-drop.git .
```

### Install

```bash
go get github.com/abhinav-ramkumar/cloud-drop
```

### Build

```bash
go build ./main.go
```

### Run

```bash
# --bucket      the S3 bucket to push your files to
# --profile     the name of your aws shared profile
# --cleanup     If true then empty folders in the uploads directory will be deleted
go run ./main.go --bucket=<bucket name> --profile=<aws-profile> --cleanup=<true/false>
```

## Milestones

- [ ] Add additional flags for advanced use cases ie multi-bucket upload config, selective uploads from uploads directory
- [ ] Set better defaults for flags
- [ ] Perform uploads within Go routines and collect status via channels
- [ ] Create an adapter that can create a client and upload to any Cloud Storage Service
