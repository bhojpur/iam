# Bhojpur IAM - Identity & Access Management

The Bhojpur IAM is used as a high performance, enterprise grade, identity &amp; access management system within the [Bhojpur.NET Platform](https://github.com/bhojpur/platform) ecosystem. It features a centralized __authentication__ and __single-sign-on__ (SSO) platform based on the OAuth 2.0 / OIDC standard protocols.

## Quick Start

You can download and run Bhojpur IAM frontend/backend services in a few minutes.

### Download

There are two methods, get the `source code` via Go subcommand `get`:

```shell
go get github.com/bhojpur/iam
```

  or `git`:

```bash
git clone https://github.com/bhojpur/iam.git
```

Finally, change the directory:

```bash
cd iam/
```

We provide two start up methods for all kinds of users.

### Manual

#### Simple Configuration

The [Bhojpur IAM](https://github.com/bhojpur/iam) requires a running relational database (e.g., MySQL server) to be operational. In this example, I have using a community edition of `MySQL` database server to be utilized by the Bhojpur IAM. You must create a database (e.g., `bhojpur`) first so that Bhojpur IAM can create required tables.

```bash
$ mysql -u root -p
mysql> CREATE DATABASE bhojpur;
Query OK, 1 row affected (0.04 sec)

mysql> quit
Bye

```

then, you need to modify the [Bhojpur IAM](https://github.com/bhojpur/iam) configuration to point to your relational database instance. Now, you should edit the `conf/app.conf` file, modify `dataSourceName` to correct relational database info, which follows this format:

```bash
username:password@tcp(database_ip:database_port)/
```

#### Run

The [Bhojpur IAM](https://github.com/bhojpur/iam) provides two run modes, the difference is in the binary image size and user prompt.

##### Development Mode

Firstly, edit the `conf/app.conf` file, set `runmode=dev`. The, build the [Bhojpur IAM](https://github.com/bhojpur/iam) web user interface front-end files:

```bash
cd pkg/webui && yarn && yarn run start

```
*The Bhojpur IAM's front-end module is built using `yarn`. You should use `yarn` instead of `npm`. It has a potential failure during building the files, if you use `npm`.*

Then, build the [Bhojpur IAM](https://github.com/bhojpur/iam) web services API back-end binary image files, change directory to project root (i.e. relative to the [Bhojpur IAM](https://github.com/bhojpur/iam)):

```bash
go run server.go web --createDatabase=true
go run client.go
```

That's it! Try to visit http://127.0.0.1:7001/. :small_airplane:  

__But make sure you always request the Bhojpur IAM backend port as 8000, when you are using SDKs.__

##### Production Mode

Firstly, edit the `conf/app.conf` file, set `runmode=prod`. Then, build the [Bhojpur IAM](https://github.com/bhojpur/iam) web user interface front-end files:

```bash
cd webui/ && yarn && yarn run build
```

then, build the [Bhojpur IAM](https://github.com/bhojpur/iam) web services API back-end binary image files, change directory to root(i.e., relative to [Bhojpur IAM](https://github.com/bhojpur/iam)):

```bash
go mod tidy
go get
go build -o bin/iamsvr server.go
go build -o bin/iamctl client.go

sudo bin/iamsvr web createDatabase=true
```

> Notice, you should visit [Bhojpur IAM](https://github.com/bhojpur/iam) web services API back-end port, default 8000. Now try to visit **http://SERVER_IP:8000/**

### Docker

The [Bhojpur IAM](https://github.com/bhojpur/iam) provides two kinds of binary images: 

- `bhojpur/platform-iam-full`, in which the [Bhojpur IAM](https://github.com/bhojpur/iam) binary, a MySQL database, and all necessary configurations are packed up. This image is for new users to have a trial on the [Bhojpur IAM](https://github.com/bhojpur/iam) platform quickly. **With this image you can start a Bhojpur IAM instance immediately with single command without any complex configuration**. **NOTE: we DO NOT recommend you to use this image in a productive environment**

- `bhojpur/platform-iam-only`: normal & graceful binary image with only [Bhojpur IAM](https://github.com/bhojpur/iam) and environment installed. 

This method requires [docker](https://docs.docker.com/get-docker/) and [docker-compose](https://docs.docker.com/compose/install/) to be installed first.

### Start Bhojpur IAM with bhojpur/platform-iam-full

if the image is not pulled, pull it from DockerHub

```shell
docker pull bhojpur/platform-iam-full
```

Start your own [Bhojpur IAM](https://github.com/bhojpur/iam) instance with

```shell
docker run -p 8000:8000 bhojpur/platform-iam-full
```

Now, you can visit http://localhost:8000 and have a try. The default account's username and password is 'admin' and '123'. Go for it!

### Start the Bhojpur IAM with bhojpur/platform-iam-only

#### modify the configurations

For the convenience of your first attempt, `docker-compose.yml` contains commands to start a database via Docker.

Thus, edit `conf/app.conf` to point to the location of your relational database(db:3306), modify `dataSourceName` to the fixed content:

```bash
dataSourceName = root:welcome1234@tcp(db:3306)/
```

> If you need to modify `conf/app.conf`, you need to re-run `docker-compose up`.

#### Run

```bash
docker-compose up
```

That's it! Try to visit http://localhost:8000/. :small_airplane:

## Detailed Documentation

We also provide complete product [documentation](https://docs.bhojpur.net/) as a reference.


## Contribute

For the [Bhojpur IAM](https://github.com/bhojpur/iam), if you have any questions, you can give Issues, or you can also directly start Pull Requests (but we recommend giving issues first to communicate with the community).

### I18N Translations

If you are contributing to the Bhojpur IAM, When you add some words in the ```pkg/webui/``` directory, please remember to add what you have added to the ```pkg/webui/src/locales/en/data.json``` file also.
